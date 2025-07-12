package service

import (
	"errors"
	"fmt"
	"go-admin/database"
	"go-admin/models"
	"go-admin/util"
	"gorm.io/gorm"
	"strconv"
)

type AuthService struct {
	db *gorm.DB
}

func NewAuthService() *AuthService {
	return &AuthService{db: database.DB}
}

func (s *AuthService) Register(data map[string]string) (*models.User, error) {
	if data["password"] != data["password_confirm"] {
		return nil, errors.New("passwords do not match")
	}

	user := &models.User{
		FirstName: data["first_name"],
		LastName:  data["last_name"],
		Email:     data["email"],
		RoleId:    1,
	}
	user.SetPassword(data["password"])

	if err := s.db.Create(user).Error; err != nil {
		return nil, fmt.Errorf("registration failed: %v", err)
	}

	return user, nil
}

func (s *AuthService) Login(data map[string]string) (*models.User, string, error) {
	var user models.User
	if err := s.db.Where("email = ?", data["email"]).First(&user).Error; err != nil {
		return nil, "", errors.New("user not found")
	}

	if err := user.ComparePassword(data["password"]); err != nil {
		return nil, "", errors.New("incorrect password")
	}

	token, err := util.GenerateJwt(strconv.Itoa(int(user.Id)))
	if err != nil {
		return nil, "", errors.New("failed to generate token")
	}

	return &user, token, nil
}

func (s *AuthService) GetUser(id string) (*models.User, error) {
	var user models.User
	userId, _ := strconv.Atoi(id)
	if err := s.db.Where("id = ?", userId).First(&user).Error; err != nil {
		return nil, errors.New("user not found")
	}

	if err := s.db.
		Preload("Role").
		Preload("Role.Permissions"). // Preload permissions
		Where("id = ?", userId).
		First(&user).Error; err != nil {
		return nil, errors.New("user not found")
	}

	return &user, nil
}

func (s *AuthService) UpdateUserInfo(id string, data map[string]string) (*models.User, error) {
	userId, _ := strconv.Atoi(id)
	user := models.User{
		Id:        uint(userId),
		FirstName: data["first_name"],
		LastName:  data["last_name"],
		Email:     data["email"],
	}

	if err := s.db.Model(&user).Updates(user).Error; err != nil {
		return nil, errors.New("update failed")
	}

	return &user, nil
}

func (s *AuthService) UpdatePassword(id string, data map[string]string) error {
	if data["password"] != data["password_confirm"] {
		return errors.New("passwords do not match")
	}

	userId, _ := strconv.Atoi(id)
	user := models.User{Id: uint(userId)}
	user.SetPassword(data["password"])

	if err := s.db.Model(&user).Updates(user).Error; err != nil {
		return errors.New("password update failed")
	}

	return nil
}