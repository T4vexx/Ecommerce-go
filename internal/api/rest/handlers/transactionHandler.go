package handlers

import (
	"encoding/json"
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
}

// MakePayment godoc
// @Summary Initiates a payment
// @Description Creates a payment request for the authenticated buyer. If an active payment exists, returns its client secret; otherwise, creates a new payment.
// @Tags Transaction
// @Produce json
// @Success 200 {object} dto.MakePaymentSuccess "Payment information including public key and client secret"
// @Failure 400 {object} dto.ErrorResponse "Error generating payment or saving payment"
// @Failure 500 {object} dto.ErrorResponse "Error generating order id"
// @Router /buyer/payment [get]
func (h *TransactionHandler) MakePayment(ctx *fiber.Ctx) error {

	user := h.svc.Auth.GetCurrentUser(ctx)
	pubKey := h.config.PubKey

	activePayment, err := h.svc.GetActivePayment(user.ID)
	if activePayment != nil && activePayment.ID > 0 {
		return ctx.Status(http.StatusOK).JSON(dto.MakePaymentSuccess{
			Message: "make payment",
			PubKey:  pubKey,
			Secret:  activePayment.ClientSecret,
		})
	}

	_, amount, err := h.userSvc.FindCart(user.ID)

	orderId, err := helper.RandomNumbers(8)
	if err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(dto.ErrorResponse{
			Message: "error generating order id",
		})
	}

	paymentResult, err := h.paymentClient.CreatePayment(amount, user.ID, orderId)
	if err != nil {
		return ctx.Status(400).JSON(dto.ErrorResponse{
			Message: "error generating payment",
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
		return ctx.Status(400).JSON(dto.ErrorResponse{
			Message: "error saving payment",
		})
	}

	return ctx.Status(http.StatusOK).JSON(dto.MakePaymentSuccess{
		Message: "create payment",
		PubKey:  pubKey,
		Secret:  paymentResult.ClientSecret,
	})
}

// VerifyPayment godoc
// @Summary Verifies payment status
// @Description Checks the status of an active payment for the authenticated buyer, updates the payment status, and creates an order if payment is successful.
// @Tags Transaction
// @Produce json
// @Success 200 {object} map[string]interface{} "Verification result along with payment response"
// @Failure 400 {object} dto.ErrorResponse "No active payment exists or error during verification"
// @Failure 500 {object} dto.ErrorResponse "Internal server error during verification"
// @Router /buyer/verify [get]
func (h *TransactionHandler) VerifyPayment(ctx *fiber.Ctx) error {
	user := h.svc.Auth.GetCurrentUser(ctx)

	activePayment, err := h.svc.GetActivePayment(user.ID)
	if err != nil || activePayment.ID == 0 {
		return ctx.Status(400).JSON(dto.ErrorResponse{
			Message: "no active payment exists",
		})
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
		return ctx.Status(fiber.StatusInternalServerError).JSON(dto.ErrorResponse{
			Message: "error verify payment",
		})
	}

	h.svc.UpdatePayment(user.ID, paymentStatus, paymentLogs)

	return ctx.Status(200).JSON(&fiber.Map{
		"message":  "create payment",
		"response": paymentRes,
	})
}
