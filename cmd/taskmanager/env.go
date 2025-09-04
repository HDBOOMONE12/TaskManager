package main

import (
	"github.com/joho/godotenv"
	"log"
	"os"
)

type Config struct {
	DatabaseURL string
}

func LoadConfig() *Config {

	_ = godotenv.Load("cmd/taskmanager/.env")
	dsn := os.Getenv("DATABASE_URL")

	if dsn == "" {
		log.Fatal("DATABASE_URL is not set")
	}

	return &Config{
		DatabaseURL: dsn,
	}
}
