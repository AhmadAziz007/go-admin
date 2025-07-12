package models

type Customer struct {
	Id    uint   `json:"id"`
	Email string `json:"email" gorm:"unique"`
	Name  string `json:"name"`
}
