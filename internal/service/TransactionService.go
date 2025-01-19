package service

import (
	"instagram-bot-live/internal/domain"
	"instagram-bot-live/internal/dto"
	"instagram-bot-live/internal/helper"
	"instagram-bot-live/internal/repository"
)

type TransactionService struct {
	Repo repository.TransactionRepository
	Auth helper.Auth
}

func (s TransactionService) GetOrders(u domain.User) ([]domain.OrderItem, error) {
	orders, err := s.Repo.FindOrders(u.ID)
	if err != nil {
		return nil, err
	}

	return orders, nil
}

func (s TransactionService) GetOrderDetails(u domain.User, id uint) (dto.SellerOrderDetails, error) {
	order, err := s.Repo.FindOrderById(u.ID, id)
	if err != nil {
		return dto.SellerOrderDetails{}, err
	}

	return order, nil
}

func (s TransactionService) GetActivePayment(uId uint) (*domain.Payment, error) {
	return s.Repo.FindInitialPayment(uId)
}

func (s TransactionService) StoreCreatePayment(input dto.CreatePaymentRequest) error {
	payment := domain.Payment{
		UserId:       input.UserId,
		Amount:       input.Amount,
		Status:       domain.PaymentStatusInitial,
		PaymentId:    input.PaymentId,
		ClientSecret: input.ClientSecret,
		OrderId:      input.OrderId,
	}

	return s.Repo.CreatePayment(&payment)
}

func (s TransactionService) UpdatePayment(userId uint, status string, paymentLog string) error {
	p, err := s.GetActivePayment(userId)
	if err != nil {
		return err
	}

	p.Status = domain.PaymentStatus(status)
	p.Response = paymentLog

	return s.Repo.UpdatePayment(p)
}

func NewTransactionService(r repository.TransactionRepository, auth helper.Auth) *TransactionService {
	return &TransactionService{
		Repo: r,
		Auth: auth,
	}
}
