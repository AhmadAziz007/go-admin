package service

import (
	"errors"
	"fmt"
	"go-admin/models"
	"gorm.io/gorm"
)

type TransactionService struct {
	db *gorm.DB
}

func NewTransactionService(db *gorm.DB) *TransactionService {
	return &TransactionService{db: db}
}

func (s *TransactionService) ValidateCartOwnership(userID uint, cartID uint) error {
	var cart models.Cart
	if err := s.db.Where("id = ? AND user_id = ?", cartID, userID).First(&cart).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errors.New("cart not found or not owned by user")
		}
		return err
	}
	return nil
}

func (s *TransactionService) SearchProduct(barcode string) (*models.Product, error) {
	var product models.Product
	result := s.db.Where("barcode = ?", barcode).First(&product)
	if result.Error != nil {
		return nil, result.Error
	}
	return &product, nil
}

func (s *TransactionService) AddToCart(userID uint, productID uint, qty float64) error {
	var product models.Product
	if err := s.db.First(&product, productID).Error; err != nil {
		return err
	}

	cart := models.Cart{
		UserID:    userID,
		ProductID: productID,
		Qty:       qty,
		Price:     product.SellPrice,
	}

	return s.db.Create(&cart).Error
}

func (s *TransactionService) DestroyCart(cartID uint) error {
	return s.db.Delete(&models.Cart{}, cartID).Error
}

func (s *TransactionService) GetCart(userID uint) ([]models.Cart, float64, error) {
	var carts []models.Cart
	result := s.db.Preload("Product").Where("user_id = ?", userID).Find(&carts)
	if result.Error != nil {
		return nil, 0, result.Error
	}

	var total float64
	for _, cart := range carts {
		total += cart.Price * cart.Qty
	}

	return carts, total, nil
}

func (s *TransactionService) PayOrder(userID, customerID uint, discountPercent, cash float64) (*models.Transaction, error) {
	carts, total, err := s.GetCart(userID)
	if err != nil {
		return nil, err
	}
	if len(carts) == 0 {
		return nil, errors.New("cart is empty")
	}

	discountAmount := (discountPercent / 100) * total
	grandTotal := total - discountAmount
	change := cash - grandTotal
	if change < 0 {
		return nil, errors.New("cash is not enough")
	}

	// Mulai transaksi
	tx := s.db.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	// Validasi stok
	for _, cart := range carts {
		var product models.Product
		if err := tx.First(&product, cart.ProductID).Error; err != nil {
			tx.Rollback()
			return nil, errors.New("product not found")
		}
		if product.Stock < cart.Qty {
			tx.Rollback()
			return nil, fmt.Errorf("stock not enough for product '%s'", product.Title)
		}
	}

	// Create Transaction
	transaction := models.Transaction{
		UserID:     userID,
		CustomerID: customerID,
		Cash:       cash,
		Change:     change,
		Discount:   discountAmount,
		GrandTotal: grandTotal,
	}

	if err := tx.Create(&transaction).Error; err != nil {
		tx.Rollback()
		return nil, err
	}

	// Create Detail dan Profit
	for _, cart := range carts {
		// Ambil data produk untuk harga beli
		var product models.Product
		if err := tx.First(&product, cart.ProductID).Error; err != nil {
			tx.Rollback()
			return nil, err
		}

		// Buat detail transaksi
		detail := models.TransactionDetail{
			TransactionID: transaction.ID,
			ProductID:     cart.ProductID,
			Qty:           cart.Qty,
			Price:         cart.Price,
		}
		if err := tx.Create(&detail).Error; err != nil {
			tx.Rollback()
			return nil, err
		}

		// Hitung profit: (sell_price - buy_price) * qty
		buyTotal := product.Price * cart.Qty
		sellTotal := product.SellPrice * cart.Qty
		profitTotal := sellTotal - buyTotal

		profit := models.Profit{
			TransactionId: transaction.ID,
			Total:         profitTotal,
		}
		if err := tx.Create(&profit).Error; err != nil {
			tx.Rollback()
			return nil, err
		}

		// Kurangi stok
		if err := tx.Model(&models.Product{}).
			Where("id = ?", cart.ProductID).
			Update("stock", gorm.Expr("stock - ?", cart.Qty)).
			Error; err != nil {
			tx.Rollback()
			return nil, err
		}
	}

	// Hapus cart
	if err := tx.Where("user_id = ?", userID).Delete(&models.Cart{}).Error; err != nil {
		tx.Rollback()
		return nil, err
	}

	if err := tx.Commit().Error; err != nil {
		return nil, err
	}

	// Ambil data lengkap
	var fullTransaction models.Transaction
	if err := s.db.Preload("TransactionDetails").First(&fullTransaction, transaction.ID).Error; err != nil {
		return nil, err
	}

	fullTransaction.DiscountPercent = discountPercent

	return &fullTransaction, nil
}
