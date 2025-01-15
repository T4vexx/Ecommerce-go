package service

import (
	"errors"
	"instagram-bot-live/config"
	"instagram-bot-live/internal/domain"
	"instagram-bot-live/internal/dto"
	"instagram-bot-live/internal/helper"
	"instagram-bot-live/internal/repository"
)

type CatalogService struct {
	Repo   repository.CatalogRepository
	Auth   helper.Auth
	Config config.AppConfig
}

func (s CatalogService) CreateCategory(input dto.CreateCategoryResquest) error {
	err := s.Repo.CreateCategory(&domain.Category{
		Name:         input.Name,
		ImageUrl:     input.ImageUrl,
		DisplayOrder: input.DisplayOrder,
		ParentId:     input.ParentId,
	})

	return err
}

func (s CatalogService) EditCategory(id int, input dto.CreateCategoryResquest) (*domain.Category, error) {
	category, err := s.Repo.FindCategoryById(id)
	if err != nil {
		return nil, errors.New("category doesn't exist")
	}

	if len(input.Name) > 0 {
		category.Name = input.Name
	}
	if len(input.ImageUrl) > 0 {
		category.ImageUrl = input.ImageUrl
	}
	if input.ParentId > 0 {
		category.ParentId = input.ParentId
	}
	if input.DisplayOrder > 0 {
		category.DisplayOrder = input.DisplayOrder
	}

	updatedCat, err := s.Repo.EditCategory(category)
	if err != nil {
		return nil, errors.New("error editing category")
	}

	return updatedCat, nil
}

func (s CatalogService) DeleteCategory(id int) error {
	err := s.Repo.DeleteCategory(id)
	if err != nil {
		return errors.New("category doesn't exist")
	}

	return nil
}

func (s CatalogService) GetCategories() ([]*domain.Category, error) {
	categories, err := s.Repo.FindCategories()
	if err != nil {
		return nil, errors.New("categories does not exist")
	}

	return categories, nil
}

func (s CatalogService) GetCategoryById(id int) (*domain.Category, error) {
	cat, err := s.Repo.FindCategoryById(id)
	if err != nil {
		return nil, errors.New("category does not exist")
	}

	return cat, nil
}

func (s CatalogService) CreateProduct(input dto.CreateProductRequest, user domain.User) error {
	err := s.Repo.CreateProduct(&domain.Product{
		Name:        input.Name,
		Description: input.Description,
		ImageUrl:    input.ImageUrl,
		Price:       input.Price,
		CategoryId:  input.CategoryId,
		UserId:      user.ID,
		Stock:       uint(input.Stock),
	})

	return err
}

func (s CatalogService) EditProduct(id int, input dto.CreateProductRequest, user domain.User) (*domain.Product, error) {
	product, err := s.Repo.FindProductById(id)
	if err != nil {
		return nil, errors.New("product doesn't exist")
	}

	if product.UserId != user.ID {
		return nil, errors.New("you don't have rights to manage this product")
	}

	if len(input.Name) > 0 {
		product.Name = input.Name
	}
	if len(input.Description) > 0 {
		product.Description = input.Description
	}
	if input.Price > 0 {
		product.Price = input.Price
	}
	if input.CategoryId > 0 {
		product.CategoryId = input.CategoryId
	}
	if len(input.ImageUrl) > 0 {
		product.ImageUrl = input.ImageUrl
	}

	updatedProduct, err := s.Repo.EditProduct(product)
	return updatedProduct, err
}

func (s CatalogService) DeleteProduct(id int, user domain.User) error {
	product, err := s.Repo.FindProductById(id)
	if err != nil {
		return errors.New("product doesn't exist")
	}

	if product.UserId != user.ID {
		return errors.New("you don't have rights to manage this product")
	}

	err = s.Repo.DeleteProduct(int(product.ID))
	if err != nil {
		return errors.New("product cannot delete")
	}

	return nil
}

func (s CatalogService) GetProductById(id int) (*domain.Product, error) {
	product, err := s.Repo.FindProductById(id)
	if err != nil {
		return nil, errors.New("product doesn't exist")
	}

	return product, nil
}

func (s CatalogService) GetProducts() ([]*domain.Product, error) {
	products, err := s.Repo.FindProducts()
	if err != nil {
		return nil, errors.New("products doesn't exist")
	}
	return products, nil
}

func (s CatalogService) GetSellerProducts(id int) ([]*domain.Product, error) {
	products, err := s.Repo.FindSellerProducts(id)
	if err != nil {
		return nil, errors.New("products doesn't exist")
	}

	return products, nil
}

func (s CatalogService) UpdateProductStock(e domain.Product) (*domain.Product, error) {
	product, err := s.Repo.FindProductById(int(e.ID))
	if err != nil {
		return nil, errors.New("product doesn't exist")
	}

	if product.UserId != e.UserId {
		return nil, errors.New("you don't have rights to manage this product")
	}

	product.Stock = uint(e.Stock)
	editProduct, err := s.Repo.EditProduct(product)
	if err != nil {
		return nil, errors.New("error updating product")
	}

	return editProduct, nil
}
