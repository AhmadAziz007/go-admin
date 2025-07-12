package middlewares

import (
	"errors"
	"fmt"
	"github.com/gofiber/fiber/v2"
	"go-admin/database"
	"go-admin/models"
	"go-admin/util"
	"strconv"
)

func IsAuthorized(c *fiber.Ctx, page string) error {
	cookie := c.Cookies("jwt")

	Id, err := util.ParseJwt(cookie)
	if err != nil {
		return err
	}

	userId, _ := strconv.Atoi(Id)

	user := models.User{
		Id: uint(userId),
	}

	database.DB.Preload("Role").Find(&user)

	role := models.Role{
		Id: user.RoleId,
	}

	database.DB.Preload("Permissions").Find(&role)

	// Debugging: Cetak ID pengguna, peran, dan izin
	fmt.Printf("User ID: %d\n", user.Id)
	fmt.Printf("Role ID: %d, Name: %s\n", role.Id, role.Name)
	fmt.Println("Permissions:")
	for _, p := range role.Permissions {
		fmt.Printf("- %s\n", p.Name)
	}

	// Periksa izin berdasarkan metode HTTP dan halaman
	requiredPermission := ""
	if c.Method() == "GET" {
		requiredPermission = "view_" + page
	} else {
		requiredPermission = "edit_" + page
	}

	// Periksa apakah pengguna memiliki izin yang diperlukan
	hasPermission := false
	for _, permission := range role.Permissions {
		if permission.Name == requiredPermission {
			hasPermission = true
			break
		}
	}

	if !hasPermission {
		errorMsg := fmt.Sprintf("unauthorized: required permission '%s'", requiredPermission)
		c.Status(fiber.StatusUnauthorized)
		return errors.New(errorMsg)
	}

	return nil
}
