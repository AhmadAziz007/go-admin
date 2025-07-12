package controller

import (
	"github.com/gofiber/fiber/v2"
	"go-admin/database"
	"go-admin/models"
)

func AllPermissions(c *fiber.Ctx) error {
	var permissions []models.Permission

	database.DB.Find(&permissions)

	return c.JSON(permissions)
}