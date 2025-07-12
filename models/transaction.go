package models

import (
	"fmt"
	"gorm.io/gorm"
	"time"
)

type Transaction struct {
	ID                 uint                `gorm:"primaryKey" json:"id"`
	UserID             uint                `json:"user_id"`
	CustomerID         uint                `json:"customer_id"`
	Invoice            string              `json:"invoice"`
	Cash               float64             `json:"cash"`
	Change             float64             `json:"change"`
	Discount           float64             `json:"discount"`
	DiscountPercent    float64             `gorm:"-" json:"discount_percent"`
	GrandTotal         float64             `json:"grand_total"`
	CreatedAt          time.Time           `json:"created_at"`
	UpdatedAt          time.Time           `json:"updated_at"`
	TransactionDetails []TransactionDetail `gorm:"foreignKey:TransactionID" json:"transaction_details"`
}

func (transaction *Transaction) AfterCreate(tx *gorm.DB) (err error) {
	year2Digit := transaction.CreatedAt.Year() % 100
	invoice := fmt.Sprintf("%d.%d.INV/ORD/%d", year2Digit, transaction.CreatedAt.Month(), transaction.ID)
	return tx.Model(transaction).Update("invoice", invoice).Error
}

type TransactionDetail struct {
	ID            uint    `gorm:"primaryKey" json:"id"`
	TransactionID uint    `json:"transaction_id"`
	ProductID     uint    `json:"product_id"`
	Qty           float64 `json:"qty"`
	Price         float64 `json:"price"`
}

func (TransactionDetail) TableName() string {
	return "transaction_details"
}
