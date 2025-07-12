package service

import (
	"errors"
	"github.com/gofiber/fiber/v2"
	"go-admin/models"
	"gorm.io/gorm"
	"math"
)

type CustomerService struct {
	db *gorm.DB
}

func NewCustomerService(db *gorm.DB) *CustomerService {
	return &CustomerService{db: db}
}

func (s *CustomerService) DropdownCustomers() ([]models.Customer, error) {
	var customers []models.Customer
	result := s.db.Find(&customers)
	return customers, result.Error
}

func (s *CustomerService) AllCustomers(page int) fiber.Map {
	limit := 5
	offset := (page - 1) * limit
	var customers []models.Customer
	s.db.Offset(offset).Limit(limit).Find(&customers)

	var total int64
	s.db.Model(&models.Customer{}).Count(&total)

	lastPage := math.Ceil(float64(total) / float64(limit))
	return fiber.Map{
		"data": customers,
		"meta": fiber.Map{
			"total":     total,
			"page":      page,
			"last_page": lastPage,
		},
	}
}

func (s *CustomerService) CreateCustomer(customer *models.Customer) error {
	result := s.db.Create(customer)
	return result.Error
}

func (s *CustomerService) GetCustomer(id uint) (models.Customer, error) {
	var customer models.Customer
	result := s.db.First(&customer, id)
	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return models.Customer{}, errors.New("customer not found")
	}
	return customer, result.Error
}

func (s *CustomerService) UpdateCustomer(id uint, updatedCustomer *models.Customer) error {
	var customer models.Customer
	if err := s.db.First(&customer, id).Error; err != nil {
		return err
	}
	return s.db.Model(&customer).Updates(updatedCustomer).Error
}

func (s *CustomerService) DeleteCustomer(id uint) error {
	result := s.db.Delete(&models.Customer{}, id)
	if result.RowsAffected == 0 {
		return errors.New("customer not found")
	}
	return result.Error
}
