package models

import "time"

type Profit struct {
	ID            uint         `gorm:"primaryKey" json:"id"`
	TransactionId uint         `json:"transaction_id"`
	Total         float64      `json:"total"`
	CreatedAt     time.Time    `json:"created_at"`
	UpdatedAt     time.Time    `json:"updated_at"`
	Transaction   *Transaction `json:"transaction" gorm:"foreignKey:TransactionId"`
}
