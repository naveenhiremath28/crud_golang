package containers

import (
	"log"
	"practise/go_fiber/internal/config"
	"practise/go_fiber/internal/database"
	"practise/go_fiber/internal/models"
	"practise/go_fiber/internal/routes"
	"practise/go_fiber/internal/service"

	"github.com/gofiber/fiber/v2"
	"go.uber.org/dig"
	"gorm.io/gorm"
)

type Container struct {
	*dig.Container
}

func NewContainer() (*Container, error) {
	c := dig.New()

	providers := []interface{}{
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

func ProvideConfig() (*config.Config, error) {
	return config.Load()
}

func ProvideDatabase(cfg *config.Config) (*gorm.DB, error) {
	db, err := database.Connect(cfg)
	if err != nil {
		log.Fatal("error while connecting to database: ", err)
	}
	db.AutoMigrate(&models.Employees{})
	return db, nil
}

func ProvideApp() (*fiber.App, error) {
	app := fiber.New()
	return app, nil
}

func ProvideRouter(app *fiber.App, cfg *config.Config, db *gorm.DB) *routes.Router {
	router := routes.NewRouter(app, cfg, db)
	router.SetupRouter()
	return router
}

func StartServer(app *fiber.App, router *routes.Router) error {
	return app.Listen(":3000")
}

func ProvideService(db *gorm.DB) *service.Service {
	return service.NewService(db)
}
