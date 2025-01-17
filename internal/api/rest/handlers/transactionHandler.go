package handlers

import (
	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
	"instagram-bot-live/internal/api/rest"
	"instagram-bot-live/internal/helper"
	"instagram-bot-live/internal/repository"
	"instagram-bot-live/internal/service"
	"instagram-bot-live/pkg/payment"
)

type TransactionHandler struct {
	svc           service.TransactionService
	userSvc       service.UserService
	paymentClient payment.PaymentClient
}

func initializeTransactionService(db *gorm.DB, auth helper.Auth) service.TransactionService {
	return service.TransactionService{
		Repo: repository.NewTransactionRepository(db),
		Auth: auth,
	}
}

func SetupTransactionRoutes(as *rest.RestHandler) {

	app := as.App
	svc := initializeTransactionService(as.DB, as.Auth)
	userSvc := service.UserService{
		Repo:   repository.NewUserRepository(as.DB),
		CRepo:  repository.NewCatalogRepository(as.DB),
		Auth:   as.Auth,
		Config: as.Config,
	}

	handler := TransactionHandler{
		svc:           svc,
		userSvc:       userSvc,
		paymentClient: as.Pc,
	}

	secRoutes := app.Group("/", as.Auth.Authorize)
	secRoutes.Get("/payment", handler.MakePayment)

	sellerRoute := app.Group("/seller", as.Auth.AuthorizeSeller)
	sellerRoute.Get("/orders", handler.GetOrders)
	sellerRoute.Get("/orders/:id", handler.GetOrderDetails)
}

func (h *TransactionHandler) MakePayment(ctx *fiber.Ctx) error {

	user := h.svc.Auth.GetCurrentUser(ctx)

	// 1. call user service get cart data to aggregate the total amount and collect payment
	_, amount, err := h.userSvc.FindCart(user.ID)

	orderId, err := helper.RandomNumbers(8)
	if err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "error generating order id",
		})
	}

	// 2. Create payment session
	sessionResult, err := h.paymentClient.CreatePayment(amount, user.ID, orderId)
	if err != nil {
		return ctx.Status(400).JSON()
	}

	// 3. Store payment session in db

	return ctx.Status(200).JSON(&fiber.Map{
		"message":     "make payment",
		"result":      sessionResult,
		"payment_url": sessionResult.URL,
	})
}

func (h *TransactionHandler) GetOrders(ctx *fiber.Ctx) error {
	return ctx.Status(200).JSON("success")
}

func (h *TransactionHandler) GetOrderDetails(ctx *fiber.Ctx) error {
	return ctx.Status(200).JSON("success")
}
