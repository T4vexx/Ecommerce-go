package handlers

import (
	"github.com/gofiber/fiber/v2"
	"instagram-bot-live/internal/api/rest"
	"instagram-bot-live/internal/domain"
	"instagram-bot-live/internal/dto"
	"instagram-bot-live/internal/repository"
	"instagram-bot-live/internal/service"
	"net/http"
	"strconv"
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
	app.Get("/products", handler.GetProducts)
	app.Get("/products/:id", handler.GetProductById)
	app.Get("/categories", handler.GetCategories)
	app.Get("/categories/:id", handler.GetCategoryById)

	// private endpoints
	selRoutes := app.Group("/seller", rh.Auth.AuthorizeSeller)
	// categories
	selRoutes.Post("/categories", handler.CreateCategory)
	selRoutes.Patch("/categories/:id", handler.EditCategory)
	selRoutes.Delete("/categories/:id", handler.DeleteCategory)

	//Product
	selRoutes.Post("/products", handler.CreateProduct)
	selRoutes.Get("/products", handler.GetSellerProducts)
	selRoutes.Patch("/products/:id", handler.UpdateStock)
	selRoutes.Put("/products/:id", handler.EditProduct)
	selRoutes.Delete("/products/:id", handler.DeleteProduct)
}

func (h CatalogHandler) GetCategories(ctx *fiber.Ctx) error {
	categories, err := h.svc.GetCategories()
	if err != nil {
		return rest.ErrorMessage(ctx, 404, err)
	}

	return rest.SuccessMessage(ctx, "Categories", categories)
}

func (h CatalogHandler) GetCategoryById(ctx *fiber.Ctx) error {
	id, _ := strconv.Atoi(ctx.Params("id"))
	cat, err := h.svc.GetCategoryById(id)
	if err != nil {
		return rest.ErrorMessage(ctx, 404, err)
	}

	return rest.SuccessMessage(ctx, "Product", cat)
}

func (h CatalogHandler) CreateCategory(ctx *fiber.Ctx) error {
	req := dto.CreateCategoryResquest{}

	err := ctx.BodyParser(&req)
	if err != nil {
		return ctx.Status(http.StatusBadRequest).JSON(&fiber.Map{
			"message": "Please provide valid inputs",
		})
	}

	err = h.svc.CreateCategory(req)
	if err != nil {
		return ctx.Status(http.StatusInternalServerError).JSON(&fiber.Map{
			"message": err.Error(),
		})
	}

	return rest.SuccessMessage(ctx, "Category created successfully", nil)
}

func (h CatalogHandler) EditCategory(ctx *fiber.Ctx) error {
	id, _ := strconv.Atoi(ctx.Params("id"))
	req := dto.CreateCategoryResquest{}

	err := ctx.BodyParser(&req)
	if err != nil {
		return ctx.Status(http.StatusBadRequest).JSON(&fiber.Map{
			"message": "Please provide valid inputs",
		})
	}

	updateCategory, err := h.svc.EditCategory(id, req)
	if err != nil {
		return ctx.Status(http.StatusInternalServerError).JSON(&fiber.Map{
			"message": err.Error(),
		})
	}

	return rest.SuccessMessage(ctx, "Category updated successfully", updateCategory)
}

func (h CatalogHandler) DeleteCategory(ctx *fiber.Ctx) error {
	id, _ := strconv.Atoi(ctx.Params("id"))
	err := h.svc.DeleteCategory(id)
	if err != nil {
		return ctx.Status(http.StatusInternalServerError).JSON(&fiber.Map{
			"message": err.Error(),
		})
	}

	return rest.SuccessMessage(ctx, "Category delete successfully", nil)
}

func (h CatalogHandler) CreateProduct(ctx *fiber.Ctx) error {
	req := dto.CreateProductRequest{}
	err := ctx.BodyParser(&req)
	if err != nil {
		return ctx.Status(http.StatusBadRequest).JSON(&fiber.Map{
			"message": "Please provide valid inputs",
		})
	}

	user := h.svc.Auth.GetCurrentUser(ctx)
	err = h.svc.CreateProduct(req, user)
	if err != nil {
		return ctx.Status(http.StatusInternalServerError).JSON(&fiber.Map{
			"message": err.Error(),
		})
	}

	return rest.SuccessMessage(ctx, "Product created successfully", nil)
}

func (h CatalogHandler) EditProduct(ctx *fiber.Ctx) error {
	id, _ := strconv.Atoi(ctx.Params("id"))
	req := dto.CreateProductRequest{}

	err := ctx.BodyParser(&req)
	if err != nil {
		return ctx.Status(http.StatusBadRequest).JSON(&fiber.Map{
			"message": "Please provide valid inputs",
		})
	}

	user := h.svc.Auth.GetCurrentUser(ctx)
	product, err := h.svc.EditProduct(id, req, user)
	if err != nil {
		return ctx.Status(http.StatusInternalServerError).JSON(&fiber.Map{
			"message": err.Error(),
		})
	}

	return rest.SuccessMessage(ctx, "Product edited successfully", product)
}

func (h CatalogHandler) DeleteProduct(ctx *fiber.Ctx) error {
	id, _ := strconv.Atoi(ctx.Params("id"))

	user := h.svc.Auth.GetCurrentUser(ctx)
	err := h.svc.DeleteProduct(id, user)
	if err != nil {
		return ctx.Status(http.StatusInternalServerError).JSON(&fiber.Map{
			"message": err.Error(),
		})
	}

	return rest.SuccessMessage(ctx, "delete product successfully", nil)
}

func (h CatalogHandler) GetProducts(ctx *fiber.Ctx) error {
	products, err := h.svc.GetProducts()
	if err != nil {
		return rest.ErrorMessage(ctx, 404, err)
	}

	return rest.SuccessMessage(ctx, "Products", products)
}

func (h CatalogHandler) GetProductById(ctx *fiber.Ctx) error {
	id, _ := strconv.Atoi(ctx.Params("id"))
	product, err := h.svc.GetProductById(id)
	if err != nil {
		return ctx.Status(http.StatusBadRequest).JSON(&fiber.Map{
			"message": err.Error(),
		})
	}

	return rest.SuccessMessage(ctx, "Product", product)
}

func (h CatalogHandler) UpdateStock(ctx *fiber.Ctx) error {
	id, _ := strconv.Atoi(ctx.Params("id"))

	req := dto.UpdateStockRequest{}
	err := ctx.BodyParser(&req)
	if err != nil {
		return ctx.Status(http.StatusBadRequest).JSON(&fiber.Map{
			"message": "Please provide valid inputs",
		})
	}

	user := h.svc.Auth.GetCurrentUser(ctx)
	product, err := h.svc.UpdateProductStock(domain.Product{
		ID:     uint(id),
		Stock:  uint(req.Stock),
		UserId: user.ID,
	})

	return rest.SuccessMessage(ctx, "get product by id endpoint", product)
}

func (h CatalogHandler) GetSellerProducts(ctx *fiber.Ctx) error {
	user := h.svc.Auth.GetCurrentUser(ctx)
	products, err := h.svc.GetSellerProducts(int(user.ID))
	if err != nil {
		return rest.ErrorMessage(ctx, 404, err)
	}

	return rest.SuccessMessage(ctx, "Seller products", products)
}
