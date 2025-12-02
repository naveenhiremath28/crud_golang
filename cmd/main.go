package main

import (
	"log"
	"practise/go_fiber/internal/config"
	database "practise/go_fiber/internal/database"
	"practise/go_fiber/internal/models"
	"practise/go_fiber/internal/routes"

	"github.com/gofiber/fiber/v2"
)

func main() {
	// Load configuration
	cfg := config.Load()

	app := fiber.New()
	err := database.Connect(cfg)
	if err != nil {
		log.Fatal("error while connecting to database: ", err)
	}
	database.DB.AutoMigrate(&models.Employees{})
	routes.SetupRouter(app, cfg)

	app.Listen(":3000")
}
