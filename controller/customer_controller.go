package controller

import (
	"go-admin/models"
	"go-admin/service"
	"net/http"
	"strconv"

	"github.com/gofiber/fiber/v2"
)

type CustomerController struct {
	service *service.CustomerService
}

func NewCustomerController(service *service.CustomerService) *CustomerController {
	return &CustomerController{service: service}
}

func (c *CustomerController) DropdownCustomers(ctx *fiber.Ctx) error {
	customers, err := c.service.DropdownCustomers()
	if err != nil {
		return ctx.Status(http.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to fetch customers",
		})
	}
	return ctx.JSON(customers)
}

func (c *CustomerController) AllCustomers(ctx *fiber.Ctx) error {

	page, _ := strconv.Atoi(ctx.Query("page", "1"))
	result := c.service.AllCustomers(page)

	return ctx.JSON(fiber.Map{
		"Code":   200,
		"Status": "OK",
		"Data":   result["data"],
		"Meta":   result["meta"],
	})
}

func (c *CustomerController) CreateCustomer(ctx *fiber.Ctx) error {
	var customer models.Customer
	if err := ctx.BodyParser(&customer); err != nil {
		return ctx.Status(http.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid request body",
		})
	}

	if err := c.service.CreateCustomer(&customer); err != nil {
		return ctx.Status(http.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to create customer",
		})
	}

	return ctx.Status(http.StatusCreated).JSON(customer)
}

func (c *CustomerController) GetCustomer(ctx *fiber.Ctx) error {
	id, err := ctx.ParamsInt("id")
	if err != nil {
		return ctx.Status(http.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid customer ID",
		})
	}

	customer, err := c.service.GetCustomer(uint(id))
	if err != nil {
		return ctx.Status(http.StatusNotFound).JSON(fiber.Map{
			"error": "Customer not found",
		})
	}

	return ctx.JSON(customer)
}

func (c *CustomerController) UpdateCustomer(ctx *fiber.Ctx) error {
	id, err := ctx.ParamsInt("id")
	if err != nil {
		return ctx.Status(http.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid customer ID",
		})
	}

	var updatedCustomer models.Customer
	if err := ctx.BodyParser(&updatedCustomer); err != nil {
		return ctx.Status(http.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid request body",
		})
	}

	if err := c.service.UpdateCustomer(uint(id), &updatedCustomer); err != nil {
		return ctx.Status(http.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to update customer",
		})
	}

	return ctx.JSON(fiber.Map{"message": "Customer updated successfully"})
}

func (c *CustomerController) DeleteCustomer(ctx *fiber.Ctx) error {
	id, err := ctx.ParamsInt("id")
	if err != nil {
		return ctx.Status(http.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid customer ID",
		})
	}

	if err := c.service.DeleteCustomer(uint(id)); err != nil {
		return ctx.Status(http.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to delete customer",
		})
	}

	return ctx.JSON(fiber.Map{"message": "Customer deleted successfully"})
}
