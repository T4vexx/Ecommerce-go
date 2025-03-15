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

// GetCategories godoc
// @Summary      Obtém todas as categorias
// @Description  Retorna uma lista com todas as categorias disponíveis
// @Tags         Catalog
// @Produce      json
// @Success      200 {object} dto.GetCategoriesSuccess "Retorna todas as categorias"
// @Failure      404 {object} dto.ErrorResponse "Retorna a mensagem de erro"
// @Router       /categories [get]
func (h CatalogHandler) GetCategories(ctx *fiber.Ctx) error {
	categories, err := h.svc.GetCategories()
	if err != nil {
		return rest.ErrorMessage(ctx, 404, err)
	}

	return rest.SuccessMessage(ctx, "Categories", categories)
}

// GetCategoryById godoc
// @Summary      Obtém categoria por ID
// @Description  Retorna os detalhes de uma categoria com base no ID informado
// @Tags         Catalog
// @Produce      json
// @Param        id   path     int  true  "ID da Categoria"
// @Success      200 {object} dto.GetCategorySuccess "Retorna a categoria por id caso o usuário tenha permissão"
// @Failure      404 {object} dto.ErrorResponse "Retorna a mensagem de erro"
// @Router       /categories/{id} [get]
func (h CatalogHandler) GetCategoryById(ctx *fiber.Ctx) error {
	id, _ := strconv.Atoi(ctx.Params("id"))
	cat, err := h.svc.GetCategoryById(id)
	if err != nil {
		return rest.ErrorMessage(ctx, 404, err)
	}

	return rest.SuccessMessage(ctx, "Product", cat)
}

// CreateCategory godoc
// @Summary      Cria uma nova categoria
// @Description  Cria uma categoria utilizando os dados fornecidos no corpo da requisição
// @Tags         Catalog
// @Accept       json
// @Produce      json
// @Param        category  body     dto.CreateCategoryResquest  true  "Dados para criação da categoria"
// @Success      200       {object} map[string]interface{}
// @Failure      400       {object} map[string]interface{}
// @Failure      500       {object} map[string]interface{}
// @Router       /seller/categories [post]
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

// EditCategory godoc
// @Summary      Atualiza uma categoria
// @Description  Atualiza os dados de uma categoria existente com base no ID informado
// @Tags         Catalog
// @Accept       json
// @Produce      json
// @Param        id        path     int                           true  "ID da Categoria"
// @Param        category  body     dto.CreateCategoryResquest    true  "Dados para atualização da categoria"
// @Success      200       {object} map[string]interface{}
// @Failure      400       {object} map[string]interface{}
// @Failure      500       {object} map[string]interface{}
// @Router       /seller/categories/{id} [patch]
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

// DeleteCategory godoc
// @Summary      Exclui uma categoria
// @Description  Remove uma categoria do sistema com base no ID informado
// @Tags         Catalog
// @Produce      json
// @Param        id   path     int  true  "ID da Categoria"
// @Success      200  {object} map[string]interface{}
// @Failure      500  {object} map[string]interface{}
// @Router       /seller/categories/{id} [delete]
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

// CreateProduct godoc
// @Summary      Cria um novo produto
// @Description  Cria um produto para o vendedor autenticado com os dados fornecidos
// @Tags         Catalog
// @Accept       json
// @Produce      json
// @Param        product  body     dto.CreateProductRequest  true  "Dados para criação do produto"
// @Success      200      {object} map[string]interface{}
// @Failure      400      {object} map[string]interface{}
// @Failure      500      {object} map[string]interface{}
// @Router       /seller/products [post]
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

// EditProduct godoc
// @Summary      Edita um produto
// @Description  Atualiza os dados de um produto do vendedor autenticado com base no ID informado
// @Tags         Catalog
// @Accept       json
// @Produce      json
// @Param        id       path     int                      true  "ID do Produto"
// @Param        product  body     dto.CreateProductRequest true  "Dados atualizados do produto"
// @Success      200      {object} map[string]interface{}
// @Failure      400      {object} map[string]interface{}
// @Failure      500      {object} map[string]interface{}
// @Router       /seller/products/{id} [put]
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

// DeleteProduct godoc
// @Summary      Exclui um produto
// @Description  Remove um produto do vendedor autenticado com base no ID informado
// @Tags         Catalog
// @Produce      json
// @Param        id   path     int  true  "ID do Produto"
// @Success      200  {object} map[string]interface{}
// @Failure      500  {object} map[string]interface{}
// @Router       /seller/products/{id} [delete]
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

// GetProducts godoc
// @Summary      Obtém todos os produtos
// @Description  Retorna uma lista com todos os produtos disponíveis
// @Tags         Catalog
// @Produce      json
// @Success      200  {object} map[string]interface{}
// @Failure      404  {object} map[string]interface{}
// @Router       /products [get]
func (h CatalogHandler) GetProducts(ctx *fiber.Ctx) error {
	products, err := h.svc.GetProducts()
	if err != nil {
		return rest.ErrorMessage(ctx, 404, err)
	}

	return rest.SuccessMessage(ctx, "Products", products)
}

// GetProductById godoc
// @Summary      Obtém produto por ID
// @Description  Retorna os detalhes de um produto com base no ID informado
// @Tags         Catalog
// @Produce      json
// @Param        id   path     int  true  "ID do Produto"
// @Success      200  {object} map[string]interface{}
// @Failure      400  {object} map[string]interface{}
// @Router       /products/{id} [get]
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

// UpdateStock godoc
// @Summary      Atualiza o estoque de um produto
// @Description  Atualiza a quantidade em estoque de um produto do vendedor autenticado
// @Tags         Catalog
// @Accept       json
// @Produce      json
// @Param        id     path     int                      true  "ID do Produto"
// @Param        stock  body     dto.UpdateStockRequest   true  "Dados para atualização do estoque"
// @Success      200    {object} map[string]interface{}
// @Failure      400    {object} map[string]interface{}
// @Failure      500    {object} map[string]interface{}
// @Router       /seller/products/{id} [patch]
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

// GetSellerProducts godoc
// @Summary      Obtém os produtos do vendedor
// @Description  Retorna a lista de produtos do vendedor autenticado
// @Tags         Catalog
// @Produce      json
// @Success      200  {object} map[string]interface{}
// @Failure      404  {object} map[string]interface{}
// @Router       /seller/products [get]
func (h CatalogHandler) GetSellerProducts(ctx *fiber.Ctx) error {
	user := h.svc.Auth.GetCurrentUser(ctx)
	products, err := h.svc.GetSellerProducts(int(user.ID))
	if err != nil {
		return rest.ErrorMessage(ctx, 404, err)
	}

	return rest.SuccessMessage(ctx, "Seller products", products)
}
