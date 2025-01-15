package repository

import (
	"errors"
	"gorm.io/gorm"
	"instagram-bot-live/internal/domain"
	"log"
)

type CatalogRepository interface {
	CreateCategory(e *domain.Category) error
	FindCategories() ([]*domain.Category, error)
	FindCategoryById(id int) (*domain.Category, error)
	EditCategory(e *domain.Category) (*domain.Category, error)
	DeleteCategory(id int) error

	CreateProduct(e *domain.Product) error
	FindProducts() ([]*domain.Product, error)
	FindProductById(id int) (*domain.Product, error)
	FindSellerProducts(id int) ([]*domain.Product, error)
	EditProduct(e *domain.Product) (*domain.Product, error)
	DeleteProduct(id int) error
}

type catalogRepository struct {
	db *gorm.DB
}

func (c catalogRepository) CreateProduct(e *domain.Product) error {
	err := c.db.Model(&domain.Product{}).Create(e).Error
	if err != nil {
		log.Printf("Err: %v", err)
		return errors.New("cannot create product")
	}

	return nil
}

func (c catalogRepository) FindProducts() ([]*domain.Product, error) {
	var products []*domain.Product
	err := c.db.Find(&products).Error
	if err != nil {
		return nil, err
	}

	return products, nil
}

func (c catalogRepository) FindProductById(id int) (*domain.Product, error) {
	var product *domain.Product
	err := c.db.First(&product, id).Error
	if err != nil {
		log.Printf("Err: %v", err)
		return nil, errors.New("cannot find product")
	}

	return product, nil
}

func (c catalogRepository) FindSellerProducts(id int) ([]*domain.Product, error) {
	var products []*domain.Product
	err := c.db.Where("user_id = ?", id).Find(&products).Error
	if err != nil {
		return nil, err
	}

	return products, nil
}

func (c catalogRepository) EditProduct(e *domain.Product) (*domain.Product, error) {
	err := c.db.Save(&e).Error
	if err != nil {
		log.Printf("Err: %v", err)
		return nil, errors.New("cannot edit product")
	}

	return e, nil
}

func (c catalogRepository) DeleteProduct(id int) error {
	err := c.db.Delete(&domain.Product{}, id).Error
	if err != nil {
		return errors.New("cannot delete product")
	}

	return nil
}

func (c catalogRepository) CreateCategory(e *domain.Category) error {
	err := c.db.Create(&e).Error
	if err != nil {
		return errors.New("Create category error ")
	}

	return nil
}

func (c catalogRepository) FindCategories() ([]*domain.Category, error) {
	var categories []*domain.Category
	err := c.db.Find(&categories).Error
	if err != nil {
		return nil, errors.New("Find categories error ")
	}

	return categories, nil
}

func (c catalogRepository) FindCategoryById(id int) (*domain.Category, error) {
	var category *domain.Category

	err := c.db.First(&category, "id = ?", id).Error
	if err != nil {
		return nil, errors.New("Find category error ")
	}

	return category, nil
}

func (c catalogRepository) EditCategory(e *domain.Category) (*domain.Category, error) {
	err := c.db.Save(&e).Error

	if err != nil {
		return nil, errors.New("Edit category error ")
	}

	return e, nil
}

func (c catalogRepository) DeleteCategory(id int) error {
	err := c.db.Delete(&domain.Category{}, "id = ?", id).Error
	if err != nil {
		return errors.New("Delete category error ")
	}

	return nil
}

func NewCatalogRepository(db *gorm.DB) CatalogRepository {
	return &catalogRepository{
		db: db,
	}
}
