package repository

import (
	"errors"
	"gorm.io/gorm/clause"
	"instagram-bot-live/internal/domain"
	"log"

	"gorm.io/gorm"
)

type UserRepository interface {
	CreateUser(u domain.User) (domain.User, error)
	FindUser(email string) (domain.User, error)
	FindUserById(id uint) (domain.User, error)
	UpdateUser(id uint, u domain.User) (domain.User, error)
	CreateBankAccount(e domain.BankAccount) error

	FindCartItems(uId uint) ([]domain.Cart, error)
	FindCartItem(uid uint, pId uint) (domain.Cart, error)
	CreateCart(c domain.Cart) error
	UpdateCart(c domain.Cart) error
	DeleteCartById(id uint) error
	DeleteCartItems(id uint) error

	CreateProfile(e domain.Address) error
	UpdateProfile(e domain.Address) error
}

type userRepository struct {
	db *gorm.DB
}

func (r userRepository) CreateProfile(e domain.Address) error {
	err := r.db.Create(&e).Error
	if err != nil {
		log.Printf("Error creating profile: %v", err)
		return errors.New("Error creating profile")
	}
	return nil
}

func (r userRepository) UpdateProfile(e domain.Address) error {
	err := r.db.Where("user_id = ?", e.UserID).Updates(e).Error
	if err != nil {
		log.Printf("Error updating profile: %v", err)
		return errors.New("Error updating profile")
	}

	return nil
}

func (r userRepository) FindCartItems(uId uint) ([]domain.Cart, error) {
	var carts []domain.Cart
	err := r.db.Where("user_id = ?", uId).Find(&carts).Error
	return carts, err
}

func (r userRepository) FindCartItem(uid uint, pId uint) (domain.Cart, error) {
	cartItem := domain.Cart{}
	err := r.db.Where("user_id = ? AND product_id = ?", uid, pId).First(&cartItem).Error
	return cartItem, err
}

func (r userRepository) CreateCart(c domain.Cart) error {
	return r.db.Create(&c).Error
}

func (r userRepository) UpdateCart(c domain.Cart) error {
	var cart domain.Cart
	err := r.db.Model(&cart).Clauses(clause.Returning{}).Where("id = ?", c.ID).Updates(c).Error
	return err
}

func (r userRepository) DeleteCartById(id uint) error {
	err := r.db.Delete(&domain.Cart{}, "id = ?", id).Error
	return err
}

func (r userRepository) DeleteCartItems(id uint) error {
	err := r.db.Where("user_id = ?", id).Delete(&domain.Cart{}).Error
	return err
}

func (r userRepository) CreateUser(u domain.User) (domain.User, error) {
	err := r.db.Create(&u).Error
	if err != nil {
		log.Printf("Create user error: %v", err)
		return domain.User{}, errors.New("Failed to create user")
	}

	return u, nil
}

func (r userRepository) FindUser(email string) (domain.User, error) {
	var user domain.User

	err := r.db.Preload("Address").First(&user, "email = ?", email).Error

	if err != nil {
		return domain.User{}, errors.New("User does not exist")
	}

	return user, nil
}

func (r userRepository) FindUserById(id uint) (domain.User, error) {
	var user domain.User

	err := r.db.Preload("Address").First(&user, id).Error

	if err != nil {
		log.Printf("find user error %v", err)
		return domain.User{}, errors.New("User does not exist")
	}

	return user, nil
}

func (r userRepository) UpdateUser(id uint, u domain.User) (domain.User, error) {
	var user domain.User

	err := r.db.Model(&user).Clauses(clause.Returning{}).Where("id = ?", id).Updates(u).Error
	if err != nil {
		log.Printf("Failed to update user: %v", err)
		return domain.User{}, errors.New("Failed Update User")
	}

	return domain.User{}, nil
}

func (r userRepository) CreateBankAccount(e domain.BankAccount) error {
	return r.db.Create(&e).Error
}

func NewUserRepository(db *gorm.DB) UserRepository {
	return &userRepository{
		db: db,
	}
}
