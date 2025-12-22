package routes

import (
	"practise/go_fiber/internal/config"
	"practise/go_fiber/internal/middlewares"
	"practise/go_fiber/internal/service"

	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

type Router struct {
	App    *fiber.App
	Config *config.Config
	DB     *gorm.DB
}

func NewRouter(app *fiber.App, cfg *config.Config, db *gorm.DB) *Router {
	return &Router{App: app, Config: cfg, DB: db}
}

func (r *Router) SetupRouter() {
	api := r.App.Group("/api")
	service := service.NewService(r.DB)
	// ---- PUBLIC ENDPOINTS (No Auth) ----
	api.Post("/login", service.LoginHandler)
	api.Post("/refresh", service.RefreshHandler)

	// ---- PROTECTED ROUTES ----
	v1 := api.Group("/v1", middlewares.KeycloakAuth(r.Config.JWKSURL))

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
