package main

import (
	"instagram-bot-live/config"
	"instagram-bot-live/internal/api"
	"log"
)

func main() {
	cfg, err := config.SetupEnv()
	if err != nil {
		log.Fatalf("Config file is not loaded properly %v\n", err)
	}

	api.StartServer(cfg)
}