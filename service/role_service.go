package service

import (
	"errors"
	"fmt"
	"go-admin/models"
	"strconv"

	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

type RoleService struct {
	db *gorm.DB
}

func NewRoleService(db *gorm.DB) *RoleService {
	return &RoleService{db: db}
}

func (s *RoleService) GetAllRoles() ([]models.Role, error) {
	var roles []models.Role
	if err := s.db.Find(&roles).Error; err != nil {
		return nil, err
	}
	return roles, nil
}

func (s *RoleService) CreateRole(roleDto fiber.Map) (*models.Role, error) {
	tx := s.db.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	name, ok := roleDto["name"].(string)
	if !ok {
		tx.Rollback()
		return nil, errors.New("invalid role name")
	}

	list, ok := roleDto["permissions"].([]interface{})
	if !ok {
		tx.Rollback()
		return nil, errors.New("invalid permissions format")
	}

	permissions := make([]models.Permission, 0, len(list))
	for _, pid := range list {
		idStr := fmt.Sprintf("%v", pid)
		id, err := strconv.ParseUint(idStr, 10, 64)
		if err != nil {
			tx.Rollback()
			return nil, errors.New("invalid permission ID format: " + idStr)
		}
		permissions = append(permissions, models.Permission{Id: uint(id)})
	}

	role := models.Role{
		Name:        name,
		Permissions: permissions,
	}

	if err := tx.Create(&role).Error; err != nil {
		tx.Rollback()
		return nil, err
	}

	// Preload permissions untuk response
	if err := tx.Preload("Permissions").First(&role, role.Id).Error; err != nil {
		tx.Rollback()
		return nil, err
	}

	if err := tx.Commit().Error; err != nil {
		return nil, err
	}

	return &role, nil
}

func (s *RoleService) GetRole(id uint) (*models.Role, error) {
	var role models.Role
	if err := s.db.Preload("Permissions").First(&role, id).Error; err != nil {
		return nil, err
	}
	return &role, nil
}

func (s *RoleService) UpdateRole(id uint, roleDto fiber.Map) (*models.Role, error) {
	tx := s.db.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	var role models.Role
	if err := tx.Preload("Permissions").First(&role, id).Error; err != nil {
		tx.Rollback()
		return nil, errors.New("role not found")
	}

	name, ok := roleDto["name"].(string)
	if ok {
		role.Name = name
	}

	if err := tx.Model(&role).Updates(role).Error; err != nil {
		tx.Rollback()
		return nil, err
	}

	if permissions, ok := roleDto["permissions"].([]interface{}); ok {
		perms := make([]models.Permission, 0, len(permissions))
		for _, pid := range permissions {
			idStr := fmt.Sprintf("%v", pid)
			id, err := strconv.ParseUint(idStr, 10, 64)
			if err != nil {
				tx.Rollback()
				return nil, errors.New("invalid permission ID format: " + idStr)
			}
			perms = append(perms, models.Permission{Id: uint(id)})
		}

		// Ganti asosiasi
		if err := tx.Model(&role).Association("Permissions").Replace(perms); err != nil {
			tx.Rollback()
			return nil, err
		}
	}

	// Preload ulang permissions
	if err := tx.Preload("Permissions").First(&role, id).Error; err != nil {
		tx.Rollback()
		return nil, err
	}

	if err := tx.Commit().Error; err != nil {
		return nil, err
	}

	return &role, nil
}

func (s *RoleService) DeleteRole(id uint) error {
	return s.db.Transaction(func(tx *gorm.DB) error {
		// Cek apakah role ada
		var role models.Role
		if err := tx.First(&role, id).Error; err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return errors.New("role not found")
			}
			return err
		}

		// Hapus relasi permission terlebih dahulu
		if err := tx.Model(&role).Association("Permissions").Clear(); err != nil {
			return errors.New("failed to clear role permissions")
		}

		// Hapus role
		if err := tx.Delete(&role).Error; err != nil {
			return errors.New("failed to delete role")
		}

		return nil
	})
}
