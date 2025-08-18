package main

import (
	"log"
	database "practise/go_fiber/internal/database"
	"practise/go_fiber/internal/routes"
	"github.com/gofiber/fiber/v2"
	"practise/go_fiber/internal/models"
)

func main() {
	app := fiber.New()
	err := database.Connect()
	if err != nil {
		log.Fatal("error while connecting to database: ", err)
	}
	database.DB.AutoMigrate(&models.Employees{})
	routes.SetupRouter(app)

	app.Listen(":3000")
}