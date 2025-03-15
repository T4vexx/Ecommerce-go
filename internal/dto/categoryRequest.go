package dto

import "instagram-bot-live/internal/domain"

type CreateCategoryResquest struct {
	Name         string `json:"name"`
	ParentId     uint   `json:"parent_id"`
	ImageUrl     string `json:"image_url"`
	DisplayOrder int    `json:"display_order"`
}

type GetCategoriesSuccess struct {
	Message string             `json:"message" example:"Categories"`
	Data    []*domain.Category `json:"data"`
}

type GetCategorySuccess struct {
	Message string           `json:"message" example:"Categories"`
	Data    *domain.Category `json:"data"`
}
