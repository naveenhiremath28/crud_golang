package main

import (
	"practise/go_fiber/internal/routes"
	"github.com/gofiber/fiber/v2"
)

func main() {
	app := fiber.New()
	routes.SetupRouter(app)
	
	app.Listen(":3000")
}