package controller

import (
	"net/http"
	"strconv"

	"go-admin/service"

	"github.com/gofiber/fiber/v2"
)

type TransactionController struct {
	service *service.TransactionService
}

func NewTransactionController(service *service.TransactionService) *TransactionController {
	return &TransactionController{service: service}
}

func (c *TransactionController) getUserID(ctx *fiber.Ctx) (uint, error) {
	userIDStr, ok := ctx.Locals("userID").(string)
	if !ok {
		return 0, fiber.NewError(http.StatusUnauthorized, "user ID not found")
	}

	userID, err := strconv.ParseUint(userIDStr, 10, 32)
	if err != nil {
		return 0, fiber.NewError(http.StatusInternalServerError, "invalid user ID format")
	}

	return uint(userID), nil
}

func (c *TransactionController) SearchProduct(ctx *fiber.Ctx) error {
	barcode := ctx.Query("barcode")
	if barcode == "" {
		return ctx.Status(http.StatusBadRequest).JSON(fiber.Map{
			"message": "Barcode is required",
		})
	}

	product, err := c.service.SearchProduct(barcode)
	if err != nil {
		return ctx.Status(http.StatusNotFound).JSON(fiber.Map{
			"message": "Product not found",
		})
	}

	return ctx.JSON(fiber.Map{
		"success": true,
		"data":    product,
	})
}

func (c *TransactionController) AddToCart(ctx *fiber.Ctx) error {
	var request struct {
		ProductID uint    `json:"product_id"`
		Qty       float64 `json:"qty"`
	}

	if err := ctx.BodyParser(&request); err != nil {
		return ctx.Status(http.StatusBadRequest).JSON(fiber.Map{
			"message": "Invalid request",
		})
	}

	userID, err := c.getUserID(ctx)
	if err != nil {
		return err
	}

	if err := c.service.AddToCart(userID, request.ProductID, request.Qty); err != nil {
		return ctx.Status(http.StatusInternalServerError).JSON(fiber.Map{
			"message": "Failed to add to cart",
		})
	}

	return ctx.JSON(fiber.Map{
		"success": true,
	})
}

func (c *TransactionController) DestroyCart(ctx *fiber.Ctx) error {
	var request struct {
		CartID uint `json:"cart_id"`
	}

	if err := ctx.BodyParser(&request); err != nil {
		return ctx.Status(http.StatusBadRequest).JSON(fiber.Map{
			"message": "Invalid request",
		})
	}

	userID, err := c.getUserID(ctx)
	if err != nil {
		return err
	}

	if err := c.service.ValidateCartOwnership(userID, request.CartID); err != nil {
		return ctx.Status(http.StatusForbidden).JSON(fiber.Map{
			"message": "You don't own this cart item",
		})
	}

	if err := c.service.DestroyCart(request.CartID); err != nil {
		return ctx.Status(http.StatusInternalServerError).JSON(fiber.Map{
			"message": "Failed to remove from cart",
		})
	}

	return ctx.JSON(fiber.Map{
		"success": true,
	})
}

func (c *TransactionController) GetCart(ctx *fiber.Ctx) error {
	userID, err := c.getUserID(ctx)
	if err != nil {
		return err
	}

	carts, total, err := c.service.GetCart(userID)
	if err != nil {
		return ctx.Status(http.StatusInternalServerError).JSON(fiber.Map{
			"message": "Failed to get cart",
		})
	}

	return ctx.JSON(fiber.Map{
		"success": true,
		"data": fiber.Map{
			"carts": carts,
			"total": total,
		},
	})
}

func (c *TransactionController) PayOrder(ctx *fiber.Ctx) error {
	var request struct {
		CustomerID uint    `json:"customer_id"`
		Discount   float64 `json:"discount"`
		Cash       float64 `json:"cash"`
	}

	if err := ctx.BodyParser(&request); err != nil {
		return ctx.Status(http.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"message": "Invalid request",
			"error":   err.Error(),
		})
	}

	userID, err := c.getUserID(ctx)
	if err != nil {
		return ctx.Status(http.StatusUnauthorized).JSON(fiber.Map{
			"success": false,
			"message": "Unauthorized",
			"error":   err.Error(),
		})
	}

	transaction, err := c.service.PayOrder(userID, request.CustomerID, request.Discount, request.Cash)
	if err != nil {
		return ctx.Status(http.StatusInternalServerError).JSON(fiber.Map{
			"success": false,
			"message": "Payment failed",
			"error":   err.Error(),
		})
	}

	return ctx.JSON(fiber.Map{
		"success": true,
		"message": "Payment successful",
		"data":    transaction,
	})
}
