package service

import (
	"errors"
	"github.com/gofiber/fiber/v2"
	"go-admin/models"
	"gorm.io/gorm"
	"math"
)

type UserService struct {
	db *gorm.DB
}

func NewUserService(db *gorm.DB) *UserService {
	return &UserService{db: db}
}

func (s *UserService) GetAllUsers(page int) fiber.Map {
	limit := 5
	offset := (page - 1) * limit

	var users []models.User
	s.db.Preload("Role").Offset(offset).Limit(limit).Find(&users)

	var total int64
	s.db.Model(&models.User{}).Count(&total)

	lastPage := math.Ceil(float64(total) / float64(limit))

	return fiber.Map{
		"data": users,
		"meta": fiber.Map{
			"total":     total,
			"page":      page,
			"last_page": lastPage,
		},
	}
}

func (s *UserService) CreateUser(user *models.User) *models.User {
	user.SetPassword("1234")
	if err := s.db.Create(user).Error; err != nil {
		return nil
	}
	return user
}

func (s *UserService) GetUser(id uint) (*models.User, error) {
	var user models.User
	if err := s.db.Preload("Role").First(&user, id).Error; err != nil {
		return nil, err
	}
	return &user, nil
}

func (s *UserService) UpdateUser(id uint, userData *models.User) (*models.User, error) {
	var user models.User
	if err := s.db.First(&user, id).Error; err != nil {
		return nil, err
	}

	// Periksa apakah email sudah digunakan oleh user lain
	if userData.Email != "" && userData.Email != user.Email {
		var existingUser models.User
		if err := s.db.Where("email = ?", userData.Email).First(&existingUser).Error; err == nil {
			if existingUser.Id != id {
				return nil, errors.New("email already in use")
			}
		}
	}

	// Update field yang diizinkan
	if userData.FirstName != "" {
		user.FirstName = userData.FirstName
	}
	if userData.LastName != "" {
		user.LastName = userData.LastName
	}
	if userData.Email != "" {
		user.Email = userData.Email
	}
	if userData.RoleId != 0 {
		user.RoleId = userData.RoleId
	}

	if err := s.db.Save(&user).Error; err != nil {
		return nil, err
	}

	// Preload role setelah update
	if err := s.db.Preload("Role").First(&user, id).Error; err != nil {
		return nil, err
	}

	return &user, nil
}

func (s *UserService) DeleteUser(id uint) error {
	result := s.db.Delete(&models.User{}, id)
	if result.RowsAffected == 0 {
		return gorm.ErrRecordNotFound
	}
	return nil
}
