package main

import (
	"instagram-bot-live/config"
	_ "instagram-bot-live/docs"
	"instagram-bot-live/internal/api"
	"log"
)

// @title           Ecommerce Go API
// @version         1.0
// @description     API de e-commerce feita em Go com Fiber
// @termsOfService  http://swagger.io/terms/

// @contact.name   Otávio Augusto
// @contact.url    http://github.com/T4vexx
// @contact.email  contato@otavioteixeira.com

// @license.name  MIT
// @license.url   https://opensource.org/licenses/MIT

// @host      localhost:3333
// @BasePath  /

func main() {
	cfg, err := config.SetupEnv()
	if err != nil {
		log.Fatalf("Config file is not loaded properly %v\n", err)
	}

	api.StartServer(cfg)
}
