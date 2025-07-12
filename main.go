package main

import (
	"go-admin/database"
	"go-admin/routes"
	"go-admin/service"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
)

func main() {
	// Initialize database
	db := database.Connect()

	// Initialize MinIO
	minioService, err := service.NewMinioService(
		"localhost:9000",
		"ytapI6KbFNtCLNmZjQ8z",
		"F4sNWuUX92S5ZSz3LGDIdc2mAOQCYuiaiPGvAgIB",
		"products",
		false,
	)
	if err != nil {
		panic(err)
	}

	app := fiber.New()

	app.Use(cors.New(cors.Config{
		AllowCredentials: true,
		AllowOrigins:     "http://localhost:8080",
		AllowHeaders:     "Origin, Content-Type, Accept, Authorization",
		AllowMethods:     "GET,POST,HEAD,PUT,DELETE,PATCH",
	}))

	// Setup routes
	routes.Setup(app, db, minioService)

	app.Listen(":8000")
}