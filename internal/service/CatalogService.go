package service

import (
	"instagram-bot-live/config"
	"instagram-bot-live/internal/helper"
	"instagram-bot-live/internal/repository"
)

type CatalogService struct {
	Repo   repository.CatalogRepository
	Auth   helper.Auth
	Config config.AppConfig
}
