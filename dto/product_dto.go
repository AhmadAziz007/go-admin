package dto

import "time"

type ProductRequest struct {
	Barcode     *string `json:"barcode"`
	Title       string  `json:"title"`
	Description string  `json:"description"`
	SellPrice   float64 `json:"sell_price"`
	Price       float64 `json:"price"`
	Stock       float64 `json:"stock"`
}

type ProductResponse struct {
	ID          uint      `json:"id"`
	Barcode     string    `json:"barcode"`
	Title       string    `json:"title"`
	Description string    `json:"description"`
	Price       float64   `json:"price"`
	SellPrice   float64   `json:"sell_price"`
	Stock       float64   `json:"stock"`
	ImgUrl      string    `json:"img_url"`
	ImageData   []byte    `json:"image_data,omitempty"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}
