package routes

import (
	"practise/go_fiber/internal/config"
	"practise/go_fiber/internal/middlewares"
	"practise/go_fiber/internal/service"

	"github.com/gofiber/fiber/v2"
)

func SetupRouter(app *fiber.App, cfg *config.Config) {
	api := app.Group("/api")
	v1 := api.Group("/v1", middlewares.KeycloakAuth(cfg.JWKSURL))

	v1.Get("/", service.ServerStatus)
	v1.Get("/listEmployees", service.ListEmployees)
	v1.Post("/addEmployee", service.AddEmployee)
	v1.Get("/getEmployee/:id", service.GetEmployee)
	v1.Patch("/updatedEmployee/:id", service.UpdateEmployee)
	v1.Delete("/deleteEmployee/:id", service.DeleteEmployee)
}
