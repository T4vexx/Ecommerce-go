package api

import (
	"github.com/gofiber/fiber/v2/middleware/cors"
	"instagram-bot-live/config"
	"instagram-bot-live/internal/api/rest"
	"instagram-bot-live/internal/api/rest/handlers"
	"instagram-bot-live/internal/domain"
	"instagram-bot-live/internal/helper"
	"log"

	"github.com/gofiber/fiber/v2"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func StartServer(config config.AppConfig) {
	app := fiber.New()

	db, err := gorm.Open(postgres.Open(config.Dsn), &gorm.Config{})
	if err != nil {
		log.Fatalf("Database connection error: %v\n", err)
	}
	log.Println("Database connected")

	// Run migrations
	err = db.AutoMigrate(&domain.User{}, &domain.BankAccount{}, &domain.Category{}, &domain.Product{})
	if err != nil {
		log.Fatalf("Error an runing migration: %v\n", err)
	}
	log.Println("migrations was successful")

	c := cors.New(cors.Config{
		AllowOrigins: "http://localhost:3000",
		AllowHeaders: "Content-Type, Accept, Authorization",
		AllowMethods: "GET, POST, PUT, PATCH, DELETE, OPTIONS, HEAD",
	})
	app.Use(c)

	auth := helper.SetupAuth(config.AppSecret)

	rh := &rest.RestHandler{
		App:    app,
		DB:     db,
		Auth:   auth,
		Config: config,
	}
	setupRouter(rh)
	app.Listen(config.ServerPort)
}

func setupRouter(rh *rest.RestHandler) {
	handlers.SetupUserRoutes(rh)
	handlers.SetupCatalogRoutes(rh)
	//catalogue

}
