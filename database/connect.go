package database

import (
	"fmt"
	"go-admin/models"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var DB *gorm.DB

func Connect() *gorm.DB {
	dsn := "host=localhost user=postgres password=root123 dbname=admin_management port=5432 sslmode=disable"
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		panic("Could not connect to database: " + err.Error())
	}

	DB = db

	err = db.AutoMigrate(
		&models.User{},
		&models.Role{},
		&models.Permission{},
		&models.Product{},
		&models.Cart{},
		&models.Customer{},
		&models.Transaction{},
		&models.TransactionDetail{},
	)
	if err != nil {
		panic("Migration failed: " + err.Error())
	}

	fmt.Println("Database connected successfully")
	return db
}
