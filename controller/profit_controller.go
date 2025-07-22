package controller

import (
	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"go-admin/dto"
	"go-admin/service"
)

type ProfitController struct {
	service  *service.ProfitService
	validate *validator.Validate
}

func NewProfitController(service *service.ProfitService) *ProfitController {
	return &ProfitController{
		service:  service,
		validate: validator.New(),
	}
}

func (c *ProfitController) FilterProfit(ctx *fiber.Ctx) error {
	var req dto.FilterDateRequest

	if err := ctx.BodyParser(&req); err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "Invalid request body",
			"error":   err.Error(),
		})
	}

	if err := c.validate.Struct(req); err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "validation failed",
			"error":   err.Error(),
		})
	}

	profits, total_profit, err := c.service.FilterProfits(req.StartDate, req.EndDate)
	if err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "Failed to get sales",
			"error":   err.Error(),
		})
	}

	return ctx.JSON(fiber.Map{
		"message":      "success",
		"profits":      profits,
		"total_profit": int(total_profit),
	})
}

func (c *ProfitController) ExportExcel(ctx *fiber.Ctx) error {
	var req dto.FilterDateRequest

	if err := ctx.BodyParser(&req); err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "Invalid request body",
			"error":   err.Error(),
		})
	}

	file, err := c.service.ExportExcel(req.StartDate, req.EndDate)
	if err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "Failed to generate Excel",
			"error":   err.Error(),
		})
	}

	// Set headers
	ctx.Set("Content-Type", "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet")
	ctx.Set("Content-Disposition", "attachment; filename=profit_report.xlsx")

	// Stream file ke client
	if _, err := file.WriteTo(ctx.Response().BodyWriter()); err != nil {
		return ctx.Status(fiber.StatusInternalServerError).SendString(err.Error())
	}

	return nil
}

func (c *ProfitController) ExportPDF(ctx *fiber.Ctx) error {
	var req dto.FilterDateRequest

	if err := ctx.BodyParser(&req); err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "Invalid request body",
			"error":   err.Error(),
		})
	}

	pdfBytes, err := c.service.ExportPDF(req.StartDate, req.EndDate)
	if err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "Failed to generate PDF",
			"error":   err.Error(),
		})
	}

	// Set headers
	ctx.Set("Content-Type", "application/pdf")
	ctx.Set("Content-Disposition", "attachment; filename=sales_report.pdf")

	return ctx.Send(pdfBytes)
}
