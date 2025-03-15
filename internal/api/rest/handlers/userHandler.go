package handlers

import (
	"instagram-bot-live/internal/api/rest"
	"instagram-bot-live/internal/dto"
	"instagram-bot-live/internal/repository"
	"instagram-bot-live/internal/service"
	"net/http"
	"strconv"

	"github.com/gofiber/fiber/v2"
)

type UserHandler struct {
	// svc UserService
	svc service.UserService
}

func SetupUserRoutes(rh *rest.RestHandler) {
	app := rh.App

	// create an instance od user service & inject to handler
	svc := service.UserService{
		Repo:   repository.NewUserRepository(rh.DB),
		Auth:   rh.Auth,
		Config: rh.Config,
		CRepo:  repository.NewCatalogRepository(rh.DB),
	}
	handler := UserHandler{
		svc: svc,
	}

	// public endpoints
	pubRoutes := app.Group("/")
	pubRoutes.Post("/register", handler.Register)
	pubRoutes.Post("/login", handler.Login)

	// private endpoints
	pvtRoutes := pubRoutes.Group("/users", rh.Auth.Authorize)
	pvtRoutes.Get("/verify", handler.GetVerificationCode)
	pvtRoutes.Post("/verify", handler.Verify)

	pvtRoutes.Get("/profile", handler.GetProfile)
	pvtRoutes.Post("/profile", handler.CreateProfile)
	pvtRoutes.Patch("/profile", handler.UpdateProfile)

	pvtRoutes.Get("/cart", handler.GetCart)
	pvtRoutes.Post("/cart", handler.AddToCart)

	pvtRoutes.Get("/order", handler.GetOrders)
	pvtRoutes.Get("/order/:id", handler.GetOrderById)

	pvtRoutes.Post("/become-seller", handler.BecomeSeller)
}

// Register godoc
// @Summary      Registra um novo usuário
// @Description  Realiza o cadastro de um novo usuário com as informações fornecidas
// @Tags         Users
// @Accept       json
// @Produce      json
// @Param        user  body     dto.UserSignup  true  "Dados para cadastro do usuário"
// @Success      200   {object} dto.UserSignupResponse "Dados retornados ao cadastrar o usuário"
// @Failure      400   {object} dto.ErrorResponse        "Erro na requisição"
// @Router       /register [post]
func (h *UserHandler) Register(ctx *fiber.Ctx) error {
	user := dto.UserSignup{}
	err := ctx.BodyParser(&user)

	if err != nil {
		return ctx.Status(http.StatusBadRequest).JSON(dto.ErrorResponse{
			Message: "Please provide valid inputs",
		})
	}

	token, err := h.svc.Signup(user)
	if err != nil {
		return ctx.Status(http.StatusBadRequest).JSON(dto.ErrorResponse{
			Message: "Error on signup",
			Reason:  err.Error(),
		})
	}

	return ctx.Status(200).JSON(dto.UserSignupResponse{
		Message: "register",
		Token:   token,
	})
}

// Login godoc
// @Summary      Login do usuário
// @Description  Realiza o login do usuário utilizando email e senha
// @Tags         Users
// @Accept       json
// @Produce      json
// @Param        credentials  body     dto.UserLogin  true  "Credenciais do usuário"
// @Success      200   {object} dto.UserSignupResponse "Dados retornados ao cadastrar o usuário"
// @Failure      400   {object} dto.ErrorResponse        "Erro na requisição"
// @Router       /login [post]
func (h *UserHandler) Login(ctx *fiber.Ctx) error {
	loginInput := dto.UserLogin{}
	err := ctx.BodyParser(&loginInput)
	if err != nil {
		return ctx.Status(http.StatusBadRequest).JSON(dto.ErrorResponse{
			Message: "Please provide valid inputs",
		})
	}

	token, err := h.svc.Login(loginInput.Email, loginInput.Password)
	if err != nil {
		return ctx.Status(http.StatusBadRequest).JSON(dto.ErrorResponse{
			Message: "Please provide valid inputs",
		})
	}

	return ctx.Status(200).JSON(dto.UserSignupResponse{
		Message: "login",
		Token:   token,
	})
}

// Verify godoc
// @Summary      Verifica código de autenticação
// @Description  Verifica o código enviado pelo usuário para ativação/validação da conta
// @Tags         Users
// @Accept       json
// @Produce      json
// @Param        verification  body     dto.VerificationCodeInput  true  "Código de verificação"
// @Success      200   {object} map[string]interface{}
// @Failure      400   {object} map[string]interface{}
// @Router       /users/verify [post]
func (h *UserHandler) Verify(ctx *fiber.Ctx) error {
	user := h.svc.Auth.GetCurrentUser(ctx)

	var req dto.VerificationCodeInput
	if err := ctx.BodyParser(&req); err != nil {
		return ctx.Status(http.StatusBadRequest).JSON(&fiber.Map{
			"message": "Please provide valid inputs",
		})
	}

	err := h.svc.VerifyCode(user.ID, req.Code)
	if err != nil {
		return ctx.Status(http.StatusBadRequest).JSON(&fiber.Map{
			"message": err.Error(),
		})
	}

	return ctx.Status(200).JSON(&fiber.Map{
		"message": "verified successfully",
	})
}

// GetVerificationCode godoc
// @Summary      Obtém código de verificação
// @Description  Gera e retorna o código de verificação para o usuário autenticado
// @Tags         Users
// @Produce      json
// @Success      200   {object} map[string]interface{}
// @Failure      500   {object} map[string]interface{}
// @Router       /users/verify [get]
func (h *UserHandler) GetVerificationCode(ctx *fiber.Ctx) error {
	user := h.svc.Auth.GetCurrentUser(ctx)

	// create verification code and update to user profile in DB
	code, err := h.svc.GetVerificationCode(user)
	if err != nil {
		return ctx.Status(http.StatusInternalServerError).JSON(&fiber.Map{
			"message": "Error on get verification code",
		})
	}

	return ctx.Status(200).JSON(&fiber.Map{
		"message": "verify",
		"data":    code,
	})
}

// CreateProfile godoc
// @Summary      Cria o perfil do usuário
// @Description  Cria o perfil do usuário autenticado com as informações fornecidas
// @Tags         Users
// @Accept       json
// @Produce      json
// @Param        profile  body     dto.ProfileInput  true  "Dados para criação do perfil"
// @Success      200   {object} map[string]interface{}
// @Failure      400   {object} map[string]interface{}
// @Router       /users/profile [post]
func (h *UserHandler) CreateProfile(ctx *fiber.Ctx) error {
	user := h.svc.Auth.GetCurrentUser(ctx)

	req := dto.ProfileInput{}
	err := ctx.BodyParser(&req)
	if err != nil {
		return ctx.Status(http.StatusBadRequest).JSON(&fiber.Map{
			"message": "Please provide valid inputs",
		})
	}

	err = h.svc.CreateProfile(user.ID, req)
	if err != nil {
		return ctx.Status(http.StatusBadRequest).JSON(&fiber.Map{
			"message": "unable to create a profile",
		})
	}

	return ctx.Status(200).JSON(&fiber.Map{
		"message": "create profile",
	})
}

// GetProfile godoc
// @Summary      Obtém o perfil do usuário
// @Description  Retorna os dados do perfil do usuário autenticado
// @Tags         Users
// @Produce      json
// @Success      200   {object} map[string]interface{}
// @Failure      500   {object} map[string]interface{}
// @Router       /users/profile [get]
func (h *UserHandler) GetProfile(ctx *fiber.Ctx) error {
	user := h.svc.Auth.GetCurrentUser(ctx)

	profile, err := h.svc.GetProfile(user.ID)
	if err != nil {
		return ctx.Status(http.StatusInternalServerError).JSON(&fiber.Map{
			"message": "Error on get profile",
		})
	}

	return ctx.Status(200).JSON(&fiber.Map{
		"message": "get profile",
		"profile": profile,
	})
}

// UpdateProfile godoc
// @Summary      Atualiza o perfil do usuário
// @Description  Atualiza as informações do perfil do usuário autenticado
// @Tags         Users
// @Accept       json
// @Produce      json
// @Param        profile  body     dto.ProfileInput  true  "Dados para atualização do perfil"
// @Success      200   {object} map[string]interface{}
// @Failure      400   {object} map[string]interface{}
// @Router       /users/profile [patch]
func (h *UserHandler) UpdateProfile(ctx *fiber.Ctx) error {
	user := h.svc.Auth.GetCurrentUser(ctx)

	req := dto.ProfileInput{}
	err := ctx.BodyParser(&req)
	if err != nil {
		return ctx.Status(http.StatusBadRequest).JSON(&fiber.Map{
			"message": "Please provide valid inputs",
		})
	}

	err = h.svc.UpdateProfile(user.ID, req)
	if err != nil {
		return ctx.Status(http.StatusBadRequest).JSON(&fiber.Map{
			"message": "unable to update profile",
		})
	}

	return rest.SuccessMessage(ctx, "profile update successfully", nil)
}

// AddToCart godoc
// @Summary      Adiciona item(s) ao carrinho
// @Description  Adiciona produtos ao carrinho do usuário autenticado
// @Tags         Users
// @Accept       json
// @Produce      json
// @Param        cart  body     dto.CreateCartRequest  true  "Dados do carrinho"
// @Success      200   {object} map[string]interface{}
// @Failure      400   {object} map[string]interface{}
// @Router       /users/cart [post]
func (h *UserHandler) AddToCart(ctx *fiber.Ctx) error {
	req := dto.CreateCartRequest{}
	err := ctx.BodyParser(&req)
	if err != nil {
		return ctx.Status(http.StatusBadRequest).JSON(&fiber.Map{
			"message": "Please provide valid inputs for the cart",
		})
	}

	user := h.svc.Auth.GetCurrentUser(ctx)
	cartItems, err := h.svc.CreateCart(req, user)
	if err != nil {
		return ctx.Status(http.StatusInternalServerError).JSON(&fiber.Map{
			"message": "Error on create cart",
		})
	}

	return rest.SuccessMessage(ctx, "Cart created successfully", cartItems)
}

// GetCart godoc
// @Summary      Obtém o carrinho
// @Description  Retorna os itens do carrinho do usuário autenticado
// @Tags         Users
// @Produce      json
// @Success      200   {object} map[string]interface{}
// @Failure      500   {object} map[string]interface{}
func (h *UserHandler) GetCart(ctx *fiber.Ctx) error {
	user := h.svc.Auth.GetCurrentUser(ctx)

	cart, _, err := h.svc.FindCart(user.ID)
	if err != nil {
		return ctx.Status(http.StatusInternalServerError).JSON(&fiber.Map{
			"message": "Error on get cart",
		})
	}

	return rest.SuccessMessage(ctx, "Cart found successfully", cart)
}

// GetOrderById godoc
// @Summary      Obtém pedido por ID
// @Description  Retorna os detalhes de um pedido específico do usuário autenticado
// @Tags         Users
// @Produce      json
// @Param        id   path     int  true  "ID do pedido"
// @Success      200  {object} map[string]interface{}
// @Failure      500  {object} map[string]interface{}
// @Router       /users/order/{id} [get]
func (h *UserHandler) GetOrderById(ctx *fiber.Ctx) error {
	user := h.svc.Auth.GetCurrentUser(ctx)
	id, _ := strconv.Atoi(ctx.Params("id"))

	order, err := h.svc.GetOrderById(uint(id), user.ID)
	if err != nil {
		return ctx.Status(http.StatusInternalServerError).JSON(&fiber.Map{
			"message": err.Error(),
		})
	}

	return ctx.Status(200).JSON(&fiber.Map{
		"message": "get order",
		"order":   order,
	})
}

// GetOrders godoc
// @Summary      Obtém todos os pedidos
// @Description  Retorna a lista de pedidos do usuário autenticado
// @Tags         Users
// @Produce      json
// @Success      200  {object} map[string]interface{}
// @Failure      500  {object} map[string]interface{}
// @Router       /users/order [get]
func (h *UserHandler) GetOrders(ctx *fiber.Ctx) error {
	user := h.svc.Auth.GetCurrentUser(ctx)

	orders, err := h.svc.GetOrders(user)
	if err != nil {
		return ctx.Status(http.StatusInternalServerError).JSON(&fiber.Map{
			"message": err.Error(),
		})
	}

	return ctx.Status(200).JSON(&fiber.Map{
		"message": "get order",
		"orders":  orders,
	})
}

// BecomeSeller godoc
// @Summary      Torna o usuário em vendedor
// @Description  Atualiza o status do usuário para vendedor, realizando as devidas validações
// @Tags         Users
// @Accept       json
// @Produce      json
// @Param        seller  body     dto.SellerInput  true  "Dados para tornar-se vendedor"
// @Success      200   {object} map[string]interface{}
// @Failure      401   {object} map[string]interface{}
// @Router       /users/become-seller [post]
func (h *UserHandler) BecomeSeller(ctx *fiber.Ctx) error {

	user := h.svc.Auth.GetCurrentUser(ctx)
	req := dto.SellerInput{}
	err := ctx.BodyParser(&req)
	if err != nil {
		return ctx.Status(http.StatusBadRequest).JSON(&fiber.Map{
			"message": "Please provide valid inputs",
		})
	}

	token, err := h.svc.BecomeSeller(user.ID, req)
	if err != nil {
		return ctx.Status(http.StatusUnauthorized).JSON(&fiber.Map{
			"message": "Failed to become seller",
		})
	}

	return ctx.Status(200).JSON(&fiber.Map{
		"message": "become seller",
		"token":   token,
	})
}
