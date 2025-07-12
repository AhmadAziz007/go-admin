package models

import (
	"gorm.io/gorm"
	"time"
)

type Cart struct {
	ID        uint      `gorm:"primaryKey" json:"id"`
	UserID    uint      `json:"user_id"`
	ProductID uint      `json:"product_id"`
	Qty       float64   `json:"qty"`
	Price     float64   `json:"price"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	Product   Product   `gorm:"foreignKey:ProductID" json:"product"`
}

func (cart *Cart) Count(db *gorm.DB) int64 {
	var total int64
	db.Model(&Cart{}).Count(&total)
	return total
}

func (cart *Cart) Take(db *gorm.DB, limit int, offset int) interface{} {
	var carts []Cart
	db.Preload("Product").Offset(offset).Limit(limit).Find(&carts)
	return carts
}
