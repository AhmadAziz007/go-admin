package service

import (
	"crypto/rand"
	"errors"
	"fmt"
	"github.com/google/uuid"
	"go-admin/dto"
	"go-admin/models"
	"mime/multipart"

	"gorm.io/gorm"
)

type ProductService struct {
	db          *gorm.DB
	minioClient *MinioService
}

func NewProductService(db *gorm.DB, minioClient *MinioService) *ProductService {
	return &ProductService{
		db:          db,
		minioClient: minioClient,
	}
}

func (s *ProductService) Create(file *multipart.FileHeader, req dto.ProductRequest) (*dto.ProductResponse, error) {
	var barcode string
	if req.Barcode == nil || *req.Barcode == "" {
		barcode = generateBarcode()
	} else {
		barcode = *req.Barcode
	}

	// Upload file
	fileSrc, err := file.Open()
	if err != nil {
		return nil, err
	}
	defer fileSrc.Close()

	imgUrl, err := s.minioClient.UploadFile(fileSrc, file.Size, file.Header.Get("Content-Type"))
	if err != nil {
		return nil, err
	}

	product := models.Product{
		Barcode:     barcode,
		Title:       req.Title,
		Description: req.Description,
		Price:       req.Price,
		SellPrice:   req.SellPrice,
		Stock:       req.Stock,
		ImgUrl:      imgUrl,
	}

	result := s.db.Create(&product)
	if result.Error != nil {
		_ = s.minioClient.DeleteFile(imgUrl)
		return nil, result.Error
	}

	return s.convertToResponse(&product)
}

func generateBarcode() string {
	uuid := uuid.New()
	return uuid.String()
}

func generateRandomBarcode() string {
	b := make([]byte, 8)
	_, err := rand.Read(b)
	if err != nil {
		return uuid.New().String()
	}
	return fmt.Sprintf("%x", b)
}

func (s *ProductService) Update(id uint, file *multipart.FileHeader, req dto.ProductRequest) (*dto.ProductResponse, error) {
	// Find existing product
	var product models.Product
	if err := s.db.First(&product, id).Error; err != nil {
		return nil, errors.New("product not found")
	}

	newImgUrl := product.ImgUrl

	// Handle file update
	if file != nil {
		// Upload new file
		fileSrc, err := file.Open()
		if err != nil {
			return nil, err
		}
		defer fileSrc.Close()

		newImgUrl, err = s.minioClient.UploadFile(fileSrc, file.Size, file.Header.Get("Content-Type"))
		if err != nil {
			return nil, err
		}

		// Delete old file
		if product.ImgUrl != "" {
			if err := s.minioClient.DeleteFile(product.ImgUrl); err != nil {
				fmt.Printf("Warning: failed to delete old image: %v\n", err)
			}
		}
	}

	// Update product
	product.Title = req.Title
	product.Description = req.Description
	product.Price = req.Price
	product.SellPrice = req.SellPrice
	product.Stock = req.Stock
	product.ImgUrl = newImgUrl

	if err := s.db.Save(&product).Error; err != nil {
		return nil, err
	}

	return s.convertToResponse(&product)
}

func (s *ProductService) Delete(id uint) error {
	var product models.Product
	if err := s.db.First(&product, id).Error; err != nil {
		return errors.New("product not found")
	}

	// Hapus gambar jika ada
	if product.ImgUrl != "" {
		if err := s.minioClient.DeleteFile(product.ImgUrl); err != nil {
			return fmt.Errorf("failed to delete image: %w", err)
		}
	}

	// Hapus produk dari database
	return s.db.Delete(&product).Error
}

func (s *ProductService) GetByID(id uint) (*dto.ProductResponse, error) {
	var product models.Product
	if err := s.db.First(&product, id).Error; err != nil {
		return nil, errors.New("product not found")
	}
	return s.convertToResponse(&product)
}

func (s *ProductService) GetAll(page, limit int) ([]dto.ProductResponse, int64, error) {
	var products []models.Product
	var total int64

	// Hitung total produk dengan benar
	s.db.Model(&models.Product{}).Count(&total)

	// Hitung offset dengan benar
	offset := (page - 1) * limit

	// Ambil data dengan pagination
	result := s.db.Offset(offset).Limit(limit).Find(&products)
	if result.Error != nil {
		return nil, 0, result.Error
	}

	// Konversi ke response TANPA mengambil gambar (untuk efisiensi)
	responses := make([]dto.ProductResponse, len(products))
	for i, p := range products {
		responses[i] = dto.ProductResponse{
			ID:          p.ID,
			Barcode:     p.Barcode,
			Title:       p.Title,
			Description: p.Description,
			Price:       p.Price,
			SellPrice:   p.SellPrice,
			Stock:       p.Stock,
			ImgUrl:      p.ImgUrl,
			CreatedAt:   p.CreatedAt,
			UpdatedAt:   p.UpdatedAt,
		}
	}

	return responses, total, nil
}

//func (s *ProductService) GetAll(page, limit int) ([]dto.ProductResponse, int64, error) {
//	var products []models.Product
//	var total int64
//
//	// Hitung total produk
//	s.db.Model(&models.Product{}).Count(&total)
//
//	// Hitung offset
//	offset := (page - 1) * limit
//
//	// Ambil data dengan pagination
//	result := s.db.Offset(offset).Limit(limit).Find(&products)
//	if result.Error != nil {
//		return nil, 0, result.Error
//	}
//
//	// Konversi ke response DENGAN mengambil gambar
//	responses := make([]dto.ProductResponse, len(products))
//
//	for i, p := range products {
//		var imageData []byte
//		var err error
//
//		// Ambil gambar hanya jika ada ImgUrl
//		if p.ImgUrl != "" {
//			imageData, err = s.minioClient.GetFile(p.ImgUrl)
//			if err != nil {
//				// Tangani error tanpa menghentikan proses
//				fmt.Printf("Error getting image for product %d: %v\n", p.ID, err)
//				// Biarkan imageData nil jika error
//			}
//		}
//
//		responses[i] = dto.ProductResponse{
//			ID:          p.ID,
//			Title:       p.Title,
//			Description: p.Description,
//			Price:       p.Price,
//			ImgUrl:      p.ImgUrl,
//			ImageData:   imageData, // Sertakan data gambar
//			CreatedAt:   p.CreatedAt,
//			UpdatedAt:   p.UpdatedAt,
//		}
//	}
//
//	return responses, total, nil
//}

func (s *ProductService) convertToResponse(p *models.Product) (*dto.ProductResponse, error) {

	var imageData []byte
	var err error

	if p.ImgUrl != "" {
		imageData, err = s.minioClient.GetFile(p.ImgUrl)
		if err != nil {
			return nil, fmt.Errorf("error getting image: %v", err)
		}
	}

	return &dto.ProductResponse{
		ID:          p.ID,
		Barcode:     p.Barcode,
		Title:       p.Title,
		Description: p.Description,
		Stock:       p.Stock,
		Price:       p.Price,
		SellPrice:   p.SellPrice,
		ImgUrl:      p.ImgUrl,
		ImageData:   imageData,
		CreatedAt:   p.CreatedAt,
		UpdatedAt:   p.UpdatedAt,
	}, nil
}
