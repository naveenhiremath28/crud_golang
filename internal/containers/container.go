package containers

import (
	"practise/go_fiber/internal/config"
	"practise/go_fiber/internal/database"
	applogger "practise/go_fiber/internal/logger"
	"practise/go_fiber/internal/models"
	"practise/go_fiber/internal/routes"
	"practise/go_fiber/internal/service"

	"github.com/gofiber/fiber/v2"
	"go.uber.org/dig"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

type Container struct {
	*dig.Container
}

func NewContainer() (*Container, error) {
	c := dig.New()

	providers := []interface{}{
		ProvideLogger,
		ProvideConfig,
		ProvideDatabase,
		ProvideApp,
		ProvideRouter,
		ProvideService,
	}

	for _, provider := range providers {
		if err := c.Provide(provider); err != nil {
			return nil, err
		}
	}

	return &Container{Container: c}, nil
}

func ProvideLogger() (*zap.SugaredLogger, error) {
	return applogger.New()
}

func ProvideConfig(log *zap.SugaredLogger) (*config.Config, error) {
	return config.Load(log)
}

func ProvideDatabase(cfg *config.Config, log *zap.SugaredLogger) (*gorm.DB, error) {
	log.Debugw("Application configuration",
		"db_host", cfg.DBHost,
		"db_port", cfg.DBPort,
		"db_user", cfg.DBUser,
		"db_name", cfg.DBName,
		"app_host", cfg.AppHost,
		"app_port", cfg.AppPort,
		"jwks_url", cfg.JWKSURL,
		"vault_url", cfg.VaultURL,
	)

	db, err := database.Connect(cfg, log)
	if err != nil {
		return nil, err
	}
	db.AutoMigrate(&models.Employees{})
	log.Info("Database schema migration completed")
	return db, nil
}

func ProvideApp() (*fiber.App, error) {
	app := fiber.New()
	return app, nil
}

func ProvideRouter(app *fiber.App, cfg *config.Config, db *gorm.DB, log *zap.SugaredLogger) *routes.Router {
	router := routes.NewRouter(app, cfg, db, log)
	router.SetupRouter()
	return router
}

func ProvideService(db *gorm.DB, cfg *config.Config, log *zap.SugaredLogger) *service.Service {
	return service.NewService(db, cfg, log)
}

func StartServer(app *fiber.App, router *routes.Router, cfg *config.Config, log *zap.SugaredLogger) error {
	addr := cfg.AppHost + ":" + cfg.AppPort
	log.Infow("Starting server", "address", addr)
	return app.Listen(addr)
}
