package controller

import (
	"fmt"
	"github.com/gofiber/fiber/v2"
	"go-admin/service"
	"go-admin/util"
	"time"
)

func Register(c *fiber.Ctx) error {
	var data map[string]string
	if err := c.BodyParser(&data); err != nil {
		return c.JSON(fiber.Map{
			"Code":   200,
			"Status": "OK",
			"Data":   fiber.Map{"error": err.Error()},
		})
	}

	authService := service.NewAuthService()
	user, err := authService.Register(data)
	if err != nil {
		return c.JSON(fiber.Map{
			"Code":   200,
			"Status": "OK",
			"Data":   fiber.Map{"error": err.Error()},
		})
	}

	return c.JSON(fiber.Map{
		"Code":   200,
		"Status": "OK",
		"Data":   user,
	})
}

func Login(c *fiber.Ctx) error {
	var data map[string]string
	if err := c.BodyParser(&data); err != nil {
		return c.JSON(fiber.Map{
			"Code":   200,
			"Status": "OK",
			"Data":   fiber.Map{"error": err.Error()},
		})
	}

	authService := service.NewAuthService()
	_, token, err := authService.Login(data)
	if err != nil {
		return c.JSON(fiber.Map{
			"Code":   200,
			"Status": "OK",
			"Data":   fiber.Map{"error": err.Error()},
		})
	}

	c.Cookie(&fiber.Cookie{
		Name:     "jwt",
		Value:    token,
		Expires:  time.Now().Add(time.Hour * 24),
		HTTPOnly: true,
	})

	return c.JSON(fiber.Map{
		"Code":   200,
		"Status": "OK",
		"Data":   "success",
	})
}

func User(c *fiber.Ctx) error {
	cookie := c.Cookies("jwt")
	id, err := util.ParseJwt(cookie)
	if err != nil {
		return c.JSON(fiber.Map{
			"Code":   200,
			"Status": "OK",
			"Data":   fiber.Map{"error": "unauthorized"},
		})
	}

	authService := service.NewAuthService()
	user, err := authService.GetUser(id)
	if err != nil {
		return c.JSON(fiber.Map{
			"Code":   200,
			"Status": "OK",
			"Data":   fiber.Map{"error": err.Error()},
		})
	}

	fmt.Printf("User Role: %+v\n", user.Role)
	if user.Role.Permissions != nil {
		fmt.Printf("Permissions: %d items\n", len(user.Role.Permissions))
	}

	return c.JSON(fiber.Map{
		"Code":   200,
		"Status": "OK",
		"Data":   user,
	})
}

func Logout(c *fiber.Ctx) error {
	c.Cookie(&fiber.Cookie{
		Name:     "jwt",
		Value:    "",
		Expires:  time.Now().Add(-time.Hour),
		HTTPOnly: true,
	})

	return c.JSON(fiber.Map{
		"Code":   200,
		"Status": "OK",
		"Data":   fiber.Map{"message": "success"},
	})
}

func UpdateInfo(c *fiber.Ctx) error {
	var data map[string]string
	if err := c.BodyParser(&data); err != nil {
		return c.JSON(fiber.Map{
			"Code":   200,
			"Status": "OK",
			"Data":   fiber.Map{"error": err.Error()},
		})
	}

	cookie := c.Cookies("jwt")
	id, _ := util.ParseJwt(cookie)

	authService := service.NewAuthService()
	user, err := authService.UpdateUserInfo(id, data)
	if err != nil {
		return c.JSON(fiber.Map{
			"Code":   200,
			"Status": "OK",
			"Data":   fiber.Map{"error": err.Error()},
		})
	}

	return c.JSON(fiber.Map{
		"Code":   200,
		"Status": "OK",
		"Data":   user,
	})
}

func UpdatePassword(c *fiber.Ctx) error {
	var data map[string]string
	if err := c.BodyParser(&data); err != nil {
		return c.JSON(fiber.Map{
			"Code":   200,
			"Status": "OK",
			"Data":   fiber.Map{"error": err.Error()},
		})
	}

	cookie := c.Cookies("jwt")
	id, _ := util.ParseJwt(cookie)

	authService := service.NewAuthService()
	if err := authService.UpdatePassword(id, data); err != nil {
		return c.JSON(fiber.Map{
			"Code":   200,
			"Status": "OK",
			"Data":   fiber.Map{"error": err.Error()},
		})
	}

	return c.JSON(fiber.Map{
		"Code":   200,
		"Status": "OK",
		"Data":   fiber.Map{"message": "password updated"},
	})
}