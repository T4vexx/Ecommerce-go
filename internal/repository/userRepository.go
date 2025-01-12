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
}

type userRepository struct {
	db *gorm.DB
}

func NewUserRepository(db *gorm.DB) UserRepository {
	return &userRepository{
		db: db,
	}
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

	err := r.db.First(&user, "email = ?", email).Error

	if err != nil {
		return domain.User{}, errors.New("User does not exist")
	}

	return user, nil
}
func (r userRepository) FindUserById(id uint) (domain.User, error) {
	var user domain.User

	err := r.db.First(&user, id).Error

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
