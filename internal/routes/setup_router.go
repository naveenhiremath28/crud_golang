package routes

import (
	"practise/go_fiber/internal/config"
	"practise/go_fiber/internal/middlewares"
	"practise/go_fiber/internal/service"

	"github.com/gofiber/fiber/v2"
)

func SetupRouter(app *fiber.App, cfg *config.Config) {
	api := app.Group("/api")

	// ---- PUBLIC ENDPOINTS (No Auth) ----
	api.Post("/login", service.LoginHandler)
	api.Post("/refresh", service.RefreshHandler)

	// ---- PROTECTED ROUTES ----
	v1 := api.Group("/v1", middlewares.KeycloakAuth(cfg.JWKSURL))

	v1.Get("/", service.ServerStatus)

	// RBAC
	v1.Get("/listEmployees",
		middlewares.RoleMiddleware("user", "manager", "admin"),
		service.ListEmployees)

	v1.Post("/addEmployee",
		middlewares.RoleMiddleware("admin"),
		service.AddEmployee)

	v1.Get("/getEmployee/:id",
		middlewares.RoleMiddleware("admin", "manager"),
		service.GetEmployee)

	v1.Patch("/updateEmployee/:id",
		middlewares.RoleMiddleware("admin", "manager"),
		service.UpdateEmployee)

	v1.Delete("/deleteEmployee/:id",
		middlewares.RoleMiddleware("admin"),
		service.DeleteEmployee)
}
