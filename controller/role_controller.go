package controller

import (
	"errors"
	"github.com/gofiber/fiber/v2"
	"go-admin/service"
	"gorm.io/gorm"
	"strconv"
)

type RoleController struct {
	service *service.RoleService
}

func NewRoleController(service *service.RoleService) *RoleController {
	return &RoleController{service: service}
}

func (c *RoleController) AllRoles(ctx *fiber.Ctx) error {
	roles, err := c.service.GetAllRoles()

	if err != nil {
		return ctx.JSON(fiber.Map{
			"Code":   200,
			"Status": "OK",
			"Data":   fiber.Map{"error": err.Error()},
		})
	}

	return ctx.JSON(fiber.Map{
		"Code":   200,
		"Status": "OK",
		"Data":   roles,
	})
}

func (c *RoleController) CreateRole(ctx *fiber.Ctx) error {
	var roleDto fiber.Map
	if err := ctx.BodyParser(&roleDto); err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"Code":   400,
			"Status": "Bad Request",
			"Data":   fiber.Map{"error": err.Error()},
		})
	}

	role, err := c.service.CreateRole(roleDto)

	if err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"Code":   500,
			"Status": "Error",
			"Data":   fiber.Map{"error": err.Error()},
		})
	}

	return ctx.JSON(fiber.Map{
		"Code":   200,
		"Status": "OK",
		"Data":   role,
	})
}

func (c *RoleController) GetRole(ctx *fiber.Ctx) error {
	id, _ := strconv.Atoi(ctx.Params("id"))
	role, err := c.service.GetRole(uint(id))

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
		"Data":   role,
	})
}

func (c *RoleController) UpdateRole(ctx *fiber.Ctx) error {
	id, err := strconv.Atoi(ctx.Params("id"))
	if err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"Code":   400,
			"Status": "Bad Request",
			"Data":   fiber.Map{"error": "invalid ID"},
		})
	}

	var roleDto fiber.Map
	if err := ctx.BodyParser(&roleDto); err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"Code":   400,
			"Status": "Bad Request",
			"Data":   fiber.Map{"error": err.Error()},
		})
	}

	updatedRole, err := c.service.UpdateRole(uint(id), roleDto)

	if err != nil {
		if err.Error() == "role not found" {
			return ctx.Status(fiber.StatusNotFound).JSON(fiber.Map{
				"Code":   404,
				"Status": "Not Found",
				"Data":   fiber.Map{"error": "role not found"},
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
		"Data":   updatedRole,
	})
}

func (c *RoleController) DeleteRole(ctx *fiber.Ctx) error {
	id, _ := strconv.Atoi(ctx.Params("id"))
	err := c.service.DeleteRole(uint(id))

	if err != nil {
		if err.Error() == "role not found" {
			return ctx.Status(404).JSON(fiber.Map{
				"Code":   404,
				"Status": "Role not found",
				"Data":   nil,
			})
		}
		return ctx.Status(500).JSON(fiber.Map{
			"Code":   500,
			"Status": "Internal Server Error: " + err.Error(),
			"Data":   nil,
		})
	}

	return ctx.JSON(fiber.Map{
		"Code":   200,
		"Status": "OK",
		"Data":   nil,
	})
}
