package models

import (
	"time"
)

type Product struct {
	ID          uint      `gorm:"primaryKey" json:"id"`
	Barcode     string    `json:"barcode"`
	Title       string    `json:"title"`
	Description string    `json:"description"`
	Stock       float64   `json:"stock"`
	Price       float64   `json:"price"`
	ImgUrl      string    `json:"img_url"`
	SellPrice   float64   `json:"sell_price"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

func (Product) TableName() string {
	return "products"
}