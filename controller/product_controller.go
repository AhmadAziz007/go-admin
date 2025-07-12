package controller

import (
	"errors"
	"go-admin/dto"
	"go-admin/service"
	"mime/multipart"
	"strconv"

	"github.com/gofiber/fiber/v2"
)

type ProductController struct {
	service *service.ProductService
}

func NewProductController(service *service.ProductService) *ProductController {
	return &ProductController{service: service}
}

func (c *ProductController) Create(ctx *fiber.Ctx) error {
	// Parse form data
	form, err := ctx.MultipartForm()
	if err != nil {
		return ctx.Status(400).JSON(fiber.Map{"error": "Invalid form data"})
	}

	// Ambil file
	files := form.File["img_url"]
	if len(files) == 0 {
		return ctx.Status(400).JSON(fiber.Map{"error": "Image is required"})
	}
	file := files[0]

	// Ambil data text
	title := form.Value["title"]
	description := form.Value["description"]
	priceStr := form.Value["price"]
	sellPriceStr := form.Value["sell_price"]
	stockStr := form.Value["stock"]

	if len(title) == 0 || len(description) == 0 || len(priceStr) == 0 || len(sellPriceStr) == 0 {
		return ctx.Status(400).JSON(fiber.Map{"error": "All fields are required"})
	}

	// Konversi price ke float
	price, err := strconv.ParseFloat(priceStr[0], 64)
	if err != nil {
		return ctx.Status(400).JSON(fiber.Map{"error": "Invalid price format"})
	}

	sellPrice, err := strconv.ParseFloat(sellPriceStr[0], 64)
	if err != nil {
		return ctx.Status(400).JSON(fiber.Map{"error": "Invalid price format"})
	}

	stock, err := strconv.ParseFloat(stockStr[0], 64)
	if err != nil {
		return ctx.Status(400).JSON(fiber.Map{"error": "Invalid stock format"})
	}

	req := dto.ProductRequest{
		Title:       title[0],
		Description: description[0],
		Price:       price,
		SellPrice:   sellPrice,
		Stock:       stock,
	}

	product, err := c.service.Create(file, req)
	if err != nil {
		return ctx.Status(500).JSON(fiber.Map{"error": err.Error()})
	}

	return ctx.Status(201).JSON(product)
}

func (c *ProductController) Update(ctx *fiber.Ctx) error {
	id, err := strconv.ParseUint(ctx.Params("id"), 10, 32)
	if err != nil {
		return ctx.Status(400).JSON(fiber.Map{"error": "Invalid ID"})
	}

	// Parse form data
	form, err := ctx.MultipartForm()
	if err != nil {
		return ctx.Status(400).JSON(fiber.Map{"error": "Invalid form data"})
	}

	var file *multipart.FileHeader
	files := form.File["img_url"]
	if len(files) > 0 {
		file = files[0]
	}

	// Ambil data text
	title := form.Value["title"]
	description := form.Value["description"]
	priceStr := form.Value["price"]
	sellPriceStr := form.Value["sell_price"]
	stockStr := form.Value["stock"]

	if len(title) == 0 || len(description) == 0 || len(priceStr) == 0 || len(sellPriceStr) == 0 {
		return ctx.Status(400).JSON(fiber.Map{"error": "All fields are required"})
	}

	price, err := strconv.ParseFloat(priceStr[0], 64)
	if err != nil {
		return ctx.Status(400).JSON(fiber.Map{"error": "Invalid price format"})
	}

	sellPrice, err := strconv.ParseFloat(sellPriceStr[0], 64)
	if err != nil {
		return ctx.Status(400).JSON(fiber.Map{"error": "Invalid price format"})
	}

	stock, err := strconv.ParseFloat(stockStr[0], 64)
	if err != nil {
		return ctx.Status(400).JSON(fiber.Map{"error": "Invalid price format"})
	}

	req := dto.ProductRequest{
		Title:       title[0],
		Description: description[0],
		Price:       price,
		SellPrice:   sellPrice,
		Stock:       stock,
	}

	product, err := c.service.Update(uint(id), file, req)
	if err != nil {
		if errors.Is(err, errors.New("product not found")) {
			return ctx.Status(404).JSON(fiber.Map{"error": "Product not found"})
		}
		return ctx.Status(500).JSON(fiber.Map{"error": err.Error()})
	}

	return ctx.JSON(product)
}

func (c *ProductController) Delete(ctx *fiber.Ctx) error {
	id, err := strconv.ParseUint(ctx.Params("id"), 10, 32)
	if err != nil {
		return ctx.Status(400).JSON(fiber.Map{
			"error": "Invalid ID",
		})
	}

	if err := c.service.Delete(uint(id)); err != nil {
		if errors.Is(err, errors.New("product not found")) {
			return ctx.Status(404).JSON(fiber.Map{"error": "Product not found"})
		}
		return ctx.Status(500).JSON(fiber.Map{"error": err.Error()})
	}

	return ctx.SendStatus(204)
}

func (c *ProductController) GetByID(ctx *fiber.Ctx) error {
	id, err := strconv.ParseUint(ctx.Params("id"), 10, 32)
	if err != nil {
		return ctx.Status(400).JSON(fiber.Map{
			"error": "Invalid ID",
		})
	}

	product, err := c.service.GetByID(uint(id))
	if err != nil {
		if errors.Is(err, errors.New("product not found")) {
			return ctx.Status(404).JSON(fiber.Map{"error": "Product not found"})
		}
		return ctx.Status(500).JSON(fiber.Map{"error": err.Error()})
	}

	return ctx.JSON(product)
}

func (c *ProductController) GetAll(ctx *fiber.Ctx) error {
	page, _ := strconv.Atoi(ctx.Query("page", "1"))
	limit, _ := strconv.Atoi(ctx.Query("limit", "5"))

	products, total, err := c.service.GetAll(page, limit)
	if err != nil {
		return ctx.Status(500).JSON(fiber.Map{"error": err.Error()})
	}

	return ctx.JSON(fiber.Map{
		"data": products,
		"meta": fiber.Map{
			"total": total,
			"page":  page,
			"limit": limit,
		},
	})
}
