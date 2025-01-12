package service

import (
	"errors"
	"instagram-bot-live/config"
	"instagram-bot-live/internal/domain"
	"instagram-bot-live/internal/dto"
	"instagram-bot-live/internal/helper"
	"instagram-bot-live/internal/repository"
	"time"
)

type UserService struct {
	Repo   repository.UserRepository
	Auth   helper.Auth
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

func (s UserService) GetVerificationCode(e domain.User) (int, error) {
	if s.isVerifiedUSer(e.ID) {
		return 0, errors.New("User already verified!")
	}

	code, err := s.Auth.GenerateCode()
	if err != nil {
		return 0, err
	}

	user := domain.User{
		Expiry: time.Now().Add(30 * time.Minute),
		Code:   code,
	}

	_, err = s.Repo.UpdateUser(e.ID, user)
	if err != nil {
		return 0, err
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

func (s UserService) VerifyCode(id uint, code int) error {
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

func (s UserService) CreateProfile(id uint, input any) error {

	return nil
}

func (s UserService) GetProfile(id uint) (*domain.User, error) {

	return nil, nil
}

func (s UserService) UpdateProfile(id uint, input any) error {

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

func (s UserService) FindaCart(id uint) ([]interface{}, error) {

	return nil, nil
}

func (s UserService) CreateCart(input any, u domain.User) ([]interface{}, error) {

	return nil, nil
}

func (s UserService) CreateOrder(u domain.User) (int, error) {

	return 0, nil
}

func (s UserService) GetOrders(u domain.User) ([]interface{}, error) {

	return nil, nil
}

func (s UserService) GetOrderById(id uint, uId uint) ([]interface{}, error) {

	return nil, nil
}
