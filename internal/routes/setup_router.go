package routes

import (
	"practise/go_fiber/internal/config"
	"practise/go_fiber/internal/middlewares"
	"practise/go_fiber/internal/service"

	"github.com/gofiber/fiber/v2"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

type Router struct {
	App    *fiber.App
	Config *config.Config
	DB     *gorm.DB
	Log    *zap.SugaredLogger
}

func NewRouter(app *fiber.App, cfg *config.Config, db *gorm.DB, log *zap.SugaredLogger) *Router {
	return &Router{App: app, Config: cfg, DB: db, Log: log}
}

func (r *Router) SetupRouter() {
	api := r.App.Group("/api")
	svc := service.NewService(r.DB, r.Config, r.Log)

	// Public endpoints
	api.Post("/login", svc.LoginHandler)
	api.Post("/refresh", svc.RefreshHandler)
	r.Log.Info("Public routes registered: /api/login, /api/refresh")

	// Protected routes
	v1 := api.Group("/v1", middlewares.KeycloakAuth(r.Log, r.Config))

	v1.Get("/", svc.ServerStatus)

	// RBAC
	v1.Get("/listEmployees",
		middlewares.RoleMiddleware("user", "manager", "admin"),
		svc.ListEmployees)

	v1.Post("/addEmployee",
		middlewares.RoleMiddleware("admin"),
		svc.AddEmployee)

	v1.Get("/getEmployee/:id",
		middlewares.RoleMiddleware("admin", "manager"),
		svc.GetEmployee)

	v1.Patch("/updateEmployee/:id",
		middlewares.RoleMiddleware("admin", "manager"),
		svc.UpdateEmployee)

	v1.Delete("/deleteEmployee/:id",
		middlewares.RoleMiddleware("admin"),
		svc.DeleteEmployee)

	r.Log.Info("Protected routes registered under /api/v1")
}
