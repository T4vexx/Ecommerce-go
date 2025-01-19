package config

import (
	"errors"
	"github.com/joho/godotenv"
	"os"
)

type AppConfig struct {
	ServerPort       string
	Dsn              string
	AppSecret        string
	TwilioAccountSid string
	TwilioAuthToken  string
	TwilioFromNumber string
	StripeSecret     string
	PubKey           string
	SuccessUrl       string
	CancelUrl        string
}

func SetupEnv() (cfg AppConfig, err error) {
	if os.Getenv("APP_ENV") == "dev" {
		godotenv.Load()
	}

	httpPort := os.Getenv("HTTP_PORT")
	if httpPort == "" {
		return AppConfig{}, errors.New("Env variables not found")
	}

	Dsn := os.Getenv("DSN")
	if Dsn == "" {
		return AppConfig{}, errors.New("Env variables not found")
	}

	appSecret := os.Getenv("APP_SECRET")
	if appSecret == "" {
		return AppConfig{}, errors.New("Env variables not found")
	}

	return AppConfig{
		ServerPort:       httpPort,
		Dsn:              Dsn,
		AppSecret:        appSecret,
		TwilioAccountSid: os.Getenv("TWILIO_ACCOUNT_SID"),
		TwilioAuthToken:  os.Getenv("TWILIO_AUTH_TOKEN"),
		TwilioFromNumber: os.Getenv("TWILIO_FROM_NUMBER"),
		StripeSecret:     os.Getenv("STRIPE_SECRET"),
		PubKey:           os.Getenv("STRIPE_PUB_KEY"),
		SuccessUrl:       os.Getenv("SUCCESS_URL"),
		CancelUrl:        os.Getenv("CANCEL_URL"),
	}, nil
}
