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

func (h *UserHandler) Register(ctx *fiber.Ctx) error {
	user := dto.UserSignup{}
	err := ctx.BodyParser(&user)

	if err != nil {
		return ctx.Status(http.StatusBadRequest).JSON(&fiber.Map{
			"message": "Please provide valid inputs",
		})
	}

	token, err := h.svc.Signup(user)
	if err != nil {
		return ctx.Status(http.StatusBadRequest).JSON(&fiber.Map{
			"message": "Error on signup",
			"reason":  err.Error(),
		})
	}

	return ctx.Status(200).JSON(&fiber.Map{
		"message": "login",
		"token":   token,
	})
}

func (h *UserHandler) Login(ctx *fiber.Ctx) error {
	loginInput := dto.UserLogin{}
	err := ctx.BodyParser(&loginInput)
	if err != nil {
		return ctx.Status(http.StatusBadRequest).JSON(&fiber.Map{
			"message": "Please provide valid inputs",
		})
	}

	token, err := h.svc.Login(loginInput.Email, loginInput.Password)
	if err != nil {
		return ctx.Status(http.StatusBadRequest).JSON(&fiber.Map{
			"message": "Please provide valid inputs",
		})
	}

	return ctx.Status(200).JSON(&fiber.Map{
		"message": "login",
		"token":   token,
	})
}

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
