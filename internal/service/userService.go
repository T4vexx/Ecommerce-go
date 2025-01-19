package service

import (
	"errors"
	"instagram-bot-live/config"
	"instagram-bot-live/internal/domain"
	"instagram-bot-live/internal/dto"
	"instagram-bot-live/internal/helper"
	"instagram-bot-live/internal/repository"
	"log"
	"time"
)

type UserService struct {
	Repo   repository.UserRepository
	Auth   helper.Auth
	CRepo  repository.CatalogRepository
	Config config.AppConfig
}

func (s UserService) Signup(input dto.UserSignup) (string, error) {
	_, err := s.findUserByEmail(input.Email)
	if err == nil {
		return "", errors.New("User already exist with provided email")
	}

	hPassword, err := s.Auth.CreateHashedPassword(input.Password)
	if err != nil {
		return "", err
	}

	user, err := s.Repo.CreateUser(domain.User{
		Email:    input.Email,
		Password: hPassword,
		Phone:    input.Phone,
	})

	return s.Auth.GenerateToken(user.ID, user.Email, user.UserType)
}

func (s UserService) Login(email string, password string) (string, error) {
	user, err := s.findUserByEmail(email)
	if err != nil {
		return "", errors.New("User does not exist with provided email")
	}

	err = s.Auth.VerifyPassword(password, user.Password)
	if err != nil {
		return "", err
	}

	return s.Auth.GenerateToken(user.ID, user.Email, user.UserType)
}

func (s UserService) findUserByEmail(email string) (*domain.User, error) {
	// perform some db opratio
	// business logic
	user, err := s.Repo.FindUser(email)

	return &user, err
}

func (s UserService) isVerifiedUSer(id uint) bool {
	currentUser, err := s.Repo.FindUserById(id)

	return err == nil && currentUser.Verified
}

func (s UserService) GetVerificationCode(e domain.User) (string, error) {
	if s.isVerifiedUSer(e.ID) {
		return "", errors.New("User already verified!")
	}

	code, err := s.Auth.GenerateCode()
	if err != nil {
		return "", err
	}

	user := domain.User{
		Expiry: time.Now().Add(30 * time.Minute),
		Code:   code,
	}

	_, err = s.Repo.UpdateUser(e.ID, user)
	if err != nil {
		return "", err
	}

	//user, err = s.Repo.FindUserById(e.ID)
	//
	//// Send SMS
	//notificationClient := notification.NewNotificationClient(s.Config)
	//notificationMessage := fmt.Sprintf("Your verification code is: %v", code)
	//err = notificationClient.SendSMS(user.Phone, notificationMessage)
	//if err != nil {
	//	return 0, errors.New("Send SMS failure! Try again later.")
	//}

	return code, nil
}

func (s UserService) VerifyCode(id uint, code string) error {
	if s.isVerifiedUSer(id) {
		return errors.New("User already verified!")
	}

	user, err := s.Repo.FindUserById(id)
	if err != nil {
		return err
	}

	if user.Code != code {
		return errors.New("Verification code does not match!")
	}

	if time.Now().After(user.Expiry) {
		return errors.New("Verification code expired!")
	}

	updateUser := domain.User{
		Verified: true,
	}
	_, err = s.Repo.UpdateUser(id, updateUser)
	if err != nil {
		return errors.New("Unable to verify user!")
	}

	return nil
}

func (s UserService) CreateProfile(id uint, input dto.ProfileInput) error {
	var user domain.User
	if input.FirstName != "" {
		user.FirstName = input.FirstName
	}
	if input.LastName != "" {
		user.LastName = input.LastName
	}

	_, err := s.Repo.UpdateUser(id, user)
	if err != nil {
		return err
	}

	address := domain.Address{
		AddressLine1: input.AddressInput.AddressLine1,
		AddressLine2: input.AddressInput.AddressLine2,
		City:         input.AddressInput.City,
		Country:      input.AddressInput.Country,
		PostalCode:   input.AddressInput.PostalCode,
		UserID:       id,
	}

	err = s.Repo.CreateProfile(address)
	if err != nil {
		return err
	}

	return nil
}

func (s UserService) GetProfile(id uint) (*domain.User, error) {

	user, err := s.Repo.FindUserById(id)
	if err != nil {
		return nil, err
	}

	return &user, nil
}

func (s UserService) UpdateProfile(id uint, input dto.ProfileInput) error {

	user, err := s.Repo.FindUserById(id)
	if err != nil {
		return err
	}

	if input.FirstName != "" {
		user.FirstName = input.FirstName
	}
	if input.LastName != "" {
		user.LastName = input.LastName
	}

	_, err = s.Repo.UpdateUser(id, user)
	if err != nil {
		return err
	}

	address := domain.Address{
		AddressLine1: input.AddressInput.AddressLine1,
		AddressLine2: input.AddressInput.AddressLine2,
		City:         input.AddressInput.City,
		Country:      input.AddressInput.Country,
		PostalCode:   input.AddressInput.PostalCode,
		UserID:       id,
	}
	err = s.Repo.UpdateProfile(address)
	if err != nil {
		return err
	}

	return nil
}

func (s UserService) BecomeSeller(id uint, input dto.SellerInput) (string, error) {
	user, _ := s.Repo.FindUserById(id)
	if user.UserType == domain.SELLER {
		return "", errors.New("User already become seller!")
	}

	_, err := s.Repo.UpdateUser(id, domain.User{
		FirstName: input.FirstName,
		LastName:  input.LastName,
		Phone:     input.PhoneNumber,
		UserType:  domain.SELLER,
	})
	if err != nil {
		return "", err
	}

	account := domain.BankAccount{
		BankAccount: input.BankAccountNumber,
		SwiftCode:   input.SwiftCode,
		PaymentType: input.PaymentType,
		UserId:      id,
	}
	err = s.Repo.CreateBankAccount(account)
	if err != nil {
		return "", err
	}

	token, err := s.Auth.GenerateToken(id, user.Email, domain.SELLER)
	if err != nil {
		return "", err
	}

	return token, nil
}

func (s UserService) FindCart(id uint) ([]domain.Cart, float64, error) {

	cartItems, err := s.Repo.FindCartItems(id)
	if err != nil {
		return nil, 0, errors.New("error on finding cart items")
	}

	var totalAmount float64
	for _, item := range cartItems {
		totalAmount += item.Price * float64(item.Qty)
	}

	return cartItems, totalAmount, err
}

func (s UserService) CreateCart(input dto.CreateCartRequest, u domain.User) ([]domain.Cart, error) {
	// check if exists, yes update, no create cart
	cart, _ := s.Repo.FindCartItem(u.ID, input.ProductId)
	if cart.ID > 0 {
		if input.ProductId == 0 {
			return nil, errors.New("Product Id is required!")
		}
		if input.Qty < 1 {
			err := s.Repo.DeleteCartById(cart.ID)
			if err != nil {
				log.Printf("Error on deleting cart item %v", err)
				return nil, errors.New("error on deleting cart item")
			}
		} else {
			cart.Qty = input.Qty
			err := s.Repo.UpdateCart(cart)
			if err != nil {
				log.Printf("Error on updating cart item %v", err)
				return nil, errors.New("error on updating cart item")
			}
		}
	} else {
		product, err := s.CRepo.FindProductById(int(input.ProductId))
		if err != nil {
			return nil, errors.New("Product not found to create cart item")
		}

		err = s.Repo.CreateCart(domain.Cart{
			ProductId: input.ProductId,
			UserId:    u.ID,
			Name:      product.Name,
			Qty:       input.Qty,
			ImageUrl:  product.ImageUrl,
			Price:     product.Price,
			SellerId:  product.UserId,
		})
		if err != nil {
			return nil, errors.New("Error on creating cart item")
		}
	}
	return s.Repo.FindCartItems(u.ID)
}

func (s UserService) CreateOrder(uId uint, orderRef string, pId string, amount float64) error {

	items, _, err := s.FindCart(uId)
	if err != nil {
		return errors.New("error on find cart items")
	}
	if len(items) == 0 {
		return errors.New("no items found, cannot create order")
	}

	var orderItems []domain.OrderItem

	for _, item := range items {
		orderItems = append(orderItems, domain.OrderItem{
			ProductId: item.ProductId,
			Qty:       item.Qty,
			Price:     item.Price,
			Name:      item.Name,
			ImageUrl:  item.ImageUrl,
			SellerId:  item.SellerId,
		})
	}

	order := domain.Order{
		UserID:         uId,
		PaymentId:      pId,
		OrderRefNumber: orderRef,
		Amount:         amount,
		Items:          orderItems,
	}
	err = s.Repo.CreateOrder(order)
	if err != nil {
		return err
	}

	err = s.Repo.DeleteCartItems(uId)
	if err != nil {
		log.Printf("Error on deleting cart item %v", err)
	}

	return nil
}

func (s UserService) GetOrders(u domain.User) ([]domain.Order, error) {
	orders, err := s.Repo.FindOrders(u.ID)
	if err != nil {
		return nil, err
	}

	return orders, nil
}

func (s UserService) GetOrderById(id uint, uId uint) (domain.Order, error) {
	order, err := s.Repo.FindOrderById(id, uId)
	if err != nil {
		return order, err
	}

	return order, nil
}
