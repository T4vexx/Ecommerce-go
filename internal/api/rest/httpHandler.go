package rest

import (
	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
	"instagram-bot-live/config"
	"instagram-bot-live/internal/helper"
)

type RestHandler struct {
	App    *fiber.App
	DB     *gorm.DB
	Auth   helper.Auth
	Config config.AppConfig
}
