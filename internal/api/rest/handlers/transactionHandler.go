package handlers

import (
	"encoding/json"
	"errors"
	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
	"instagram-bot-live/config"
	"instagram-bot-live/internal/api/rest"
	"instagram-bot-live/internal/dto"
	"instagram-bot-live/internal/helper"
	"instagram-bot-live/internal/repository"
	"instagram-bot-live/internal/service"
	"instagram-bot-live/pkg/payment"
	"net/http"
)

type TransactionHandler struct {
	svc           service.TransactionService
	userSvc       service.UserService
	paymentClient payment.PaymentClient
	config        config.AppConfig
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
		config:        as.Config,
	}

	secRoutes := app.Group("/buyer", as.Auth.Authorize)
	secRoutes.Get("/payment", handler.MakePayment)
	secRoutes.Get("/verify", handler.VerifyPayment)

	//sellerRoute := app.Group("/seller", as.Auth.AuthorizeSeller)
	//sellerRoute.Get("/orders", handler.GetOrders)
	//sellerRoute.Get("/orders/:id", handler.GetOrderDetails)
}

func (h *TransactionHandler) MakePayment(ctx *fiber.Ctx) error {

	user := h.svc.Auth.GetCurrentUser(ctx)
	pubKey := h.config.PubKey

	activePayment, err := h.svc.GetActivePayment(user.ID)
	if activePayment != nil && activePayment.ID > 0 {
		return ctx.Status(http.StatusOK).JSON(&fiber.Map{
			"message": "make payment",
			"pubKey":  pubKey,
			"secret":  activePayment.ClientSecret,
		})
	}

	_, amount, err := h.userSvc.FindCart(user.ID)

	orderId, err := helper.RandomNumbers(8)
	if err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "error generating order id",
		})
	}

	paymentResult, err := h.paymentClient.CreatePayment(amount, user.ID, orderId)
	if err != nil {
		return ctx.Status(400).JSON(&fiber.Map{
			"message": "error generating payment",
		})
	}

	err = h.svc.StoreCreatePayment(dto.CreatePaymentRequest{
		UserId:       user.ID,
		Amount:       amount,
		ClientSecret: paymentResult.ClientSecret,
		PaymentId:    paymentResult.ID,
		OrderId:      orderId,
	})
	if err != nil {
		return ctx.Status(400).JSON(&fiber.Map{
			"message": "error saving payment",
		})
	}

	return ctx.Status(200).JSON(&fiber.Map{
		"message": "create payment",
		"pubKey":  pubKey,
		"secret":  paymentResult.ClientSecret,
	})
}

func (h *TransactionHandler) VerifyPayment(ctx *fiber.Ctx) error {
	user := h.svc.Auth.GetCurrentUser(ctx)

	activePayment, err := h.svc.GetActivePayment(user.ID)
	if err != nil || activePayment.ID == 0 {
		return ctx.Status(400).JSON(errors.New("no active payment exist"))
	}

	paymentRes, err := h.paymentClient.GetPaymentStatus(activePayment.PaymentId)
	paymentJson, _ := json.Marshal(paymentRes)
	paymentLogs := string(paymentJson)
	paymentStatus := "failed"

	if paymentRes.Status == "succeeded" {
		paymentStatus = "success"
		err = h.userSvc.CreateOrder(user.ID, activePayment.OrderId, activePayment.PaymentId, activePayment.Amount)
	}

	if err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(&fiber.Map{
			"message": "error verifying payment",
		})
	}

	h.svc.UpdatePayment(user.ID, paymentStatus, paymentLogs)

	return ctx.Status(200).JSON(&fiber.Map{
		"message":  "create payment",
		"response": paymentRes,
	})
}

//func (h *TransactionHandler) GetOrders(ctx *fiber.Ctx) error {
//	return ctx.Status(200).JSON("success")
//}
//
//func (h *TransactionHandler) GetOrderDetails(ctx *fiber.Ctx) error {
//	return ctx.Status(200).JSON("success")
//}
