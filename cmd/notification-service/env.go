package main

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	TelegramToken string
	Port          string
	DatabaseURL   string
	BaseUrl       string
}

func LoadConfig() *Config {
	_ = godotenv.Load("cmd/notification-service/.env")

	token := os.Getenv("TELEGRAM_BOT_TOKEN")
	if token == "" {
		log.Fatal("TELEGRAM_TOKEN is not set")
	}

	port := os.Getenv("PORT")
	if port == "" {
		port = "8081"
	}

	dsn := os.Getenv("DATABASE_URL")
	if dsn == "" {
		log.Fatal("DATABASE_URL is not set")
	}

	baseUrl := os.Getenv("TASK_SERVICE_GRPC_URL")
	if baseUrl == "" {
		log.Fatal("TASK_SERVICE_GRPC_URL is not set")
	}

	return &Config{
		TelegramToken: token,
		Port:          port,
		DatabaseURL:   dsn,
		BaseUrl:       baseUrl,
	}
}
