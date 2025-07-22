package controller

import (
	"github.com/gofiber/fiber/v2"
	"go-admin/service"
)

type DashboardController struct {
	service *service.DashboardService
}

func NewDashboardController(service *service.DashboardService) *DashboardController {
	return &DashboardController{service: service}
}

func (dc *DashboardController) GetDashboard(c *fiber.Ctx) error {
	data, err := dc.service.GetDashboardData()
	if err != nil {
		return c.Status(500).JSON(fiber.Map{
			"message": "Failed to fetch dashboard data",
			"error":   err.Error(),
		})
	}

	return c.JSON(fiber.Map{
		"message": "success",
		"data":    data,
	})
}
