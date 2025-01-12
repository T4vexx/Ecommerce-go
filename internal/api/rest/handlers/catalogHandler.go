package handlers

import (
	"github.com/gofiber/fiber/v2"
	"instagram-bot-live/internal/api/rest"
	"instagram-bot-live/internal/repository"
	"instagram-bot-live/internal/service"
)

package handlers

import (
"instagram-bot-live/internal/api/rest"
"instagram-bot-live/internal/dto"
"instagram-bot-live/internal/repository"
"instagram-bot-live/internal/service"
"net/http"

"github.com/gofiber/fiber/v2"
)

type CatalogHandler struct {
	svc service.CatalogService
}

func SetupCatalogRoutes(rh *rest.RestHandler) {
	app := rh.App

	// create an instance od user service & inject to handler
	svc := service.CatalogService{
		Repo:   repository.NewCatalogRepository(rh.DB),
		Auth:   rh.Auth,
		Config: rh.Config,
	}
	handler := CatalogHandler{
		svc: svc,
	}

	// public endpoints
	app.Get("/products")
	app.Get("/products/:id")
	app.Get("/categories")
	app.Get("/categories/:id")

	// private endpoints
	selRoutes := app.Group("/seller")
	// categories
	selRoutes.Post("/categories")
	selRoutes.Patch("/categories/:id")
	selRoutes.Delete("/categories/:id")

	//Product
	selRoutes.Post("/products")
	selRoutes.Get("/products")
	selRoutes.Get("/products/:id")
	selRoutes.Patch("/products/:id")
	selRoutes.Put("/products/:id")
	selRoutes.Delete("/products/:id")
}

func (h CatalogHandler) CreateCategories(ctx *fiber.Ctx) error {

}