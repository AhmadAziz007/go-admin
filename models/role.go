package models

type Role struct {
	Id          uint         `json:"id"`
	Name        string       `json:"name"`
	Permissions []Permission `json:"permissions,omitempty" gorm:"many2many:role_permissions;"`
}
