package controller

import (
	"errors"
	"github.com/gofiber/fiber/v2"
	"go-admin/middlewares"
	"go-admin/models"
	"go-admin/service"
	"gorm.io/gorm"
	"strconv"
)

type UserController struct {
	service *service.UserService
}

func NewUserController(service *service.UserService) *UserController {
	return &UserController{service: service}
}

func (c *UserController) AllUsers(ctx *fiber.Ctx) error {
	if err := middlewares.IsAuthorized(ctx, "users"); err != nil {
		return err
	}

	page, _ := strconv.Atoi(ctx.Query("page", "1"))
	result := c.service.GetAllUsers(page)

	return ctx.JSON(fiber.Map{
		"Code":   200,
		"Status": "OK",
		"Data":   result["data"],
		"Meta":   result["meta"],
	})
}

func (c *UserController) CreateUser(ctx *fiber.Ctx) error {
	if err := middlewares.IsAuthorized(ctx, "users"); err != nil {
		return err
	}

	var user models.User
	if err := ctx.BodyParser(&user); err != nil {
		return ctx.JSON(fiber.Map{
			"Code":   200,
			"Status": "OK",
			"Data":   fiber.Map{"error": err.Error()},
		})
	}

	createdUser := c.service.CreateUser(&user)

	return ctx.JSON(fiber.Map{
		"Code":   200,
		"Status": "OK",
		"Data":   createdUser,
	})
}

func (c *UserController) GetUser(ctx *fiber.Ctx) error {
	id, _ := strconv.Atoi(ctx.Params("id"))
	user, err := c.service.GetUser(uint(id))

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return ctx.Status(401).JSON(fiber.Map{
				"Code":   401,
				"Status": "id not found",
				"Data":   nil,
			})
		}
		return ctx.Status(500).JSON(fiber.Map{
			"Code":   500,
			"Status": "Internal Server Error",
			"Data":   nil,
		})
	}

	return ctx.JSON(fiber.Map{
		"Code":   200,
		"Status": "OK",
		"Data":   user,
	})
}

func (c *UserController) UpdateUser(ctx *fiber.Ctx) error {
	id, err := strconv.Atoi(ctx.Params("id"))
	if err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"Code":   400,
			"Status": "Bad Request",
			"Data":   fiber.Map{"error": "Invalid ID format"},
		})
	}

	// Parse body request
	var userData models.User
	if err := ctx.BodyParser(&userData); err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"Code":   400,
			"Status": "Bad Request",
			"Data":   fiber.Map{"error": err.Error()},
		})
	}

	updatedUser, err := c.service.UpdateUser(uint(id), &userData)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return ctx.Status(fiber.StatusNotFound).JSON(fiber.Map{
				"Code":   404,
				"Status": "Not Found",
				"Data":   fiber.Map{"error": "User not found"},
			})
		}

		if err.Error() == "email already in use" {
			return ctx.Status(fiber.StatusConflict).JSON(fiber.Map{
				"Code":   409,
				"Status": "Conflict",
				"Data":   fiber.Map{"error": "Email already in use"},
			})
		}

		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"Code":   500,
			"Status": "Internal Server Error",
			"Data":   fiber.Map{"error": err.Error()},
		})
	}

	return ctx.JSON(fiber.Map{
		"Code":   200,
		"Status": "OK",
		"Data":   updatedUser,
	})
}

func (c *UserController) DeleteUser(ctx *fiber.Ctx) error {
	id, _ := strconv.Atoi(ctx.Params("id"))
	err := c.service.DeleteUser(uint(id))

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return ctx.Status(401).JSON(fiber.Map{
				"Code":   401,
				"Status": "id not found",
				"Data":   nil,
			})
		}
		return ctx.Status(500).JSON(fiber.Map{
			"Code":   500,
			"Status": "Internal Server Error",
			"Data":   nil,
		})
	}

	return ctx.JSON(fiber.Map{
		"Code":   200,
		"Status": "OK",
		"Data":   nil,
	})
}
