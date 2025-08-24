package db

import (
	"database/sql"
	"log"
	"os"
	"time"
)
import "github.com/joho/godotenv"

func Init() *sql.DB {
	_ = godotenv.Load()
	dsn := os.Getenv("DATABASE_URL")

	if dsn == "" {
		log.Fatal("DATABASE_URL is not set")
	}

	db, err := sql.Open("pgx", dsn)
	if err != nil {
		log.Fatalf("cannot open DB: %v", err)
	}
	db.SetMaxOpenConns(10)
	db.SetMaxIdleConns(5)
	db.SetConnMaxLifetime(time.Minute * 30)
	err = db.Ping()

	if err := db.Ping(); err != nil {
		log.Fatalf("cannot connect to DB: %v", err)
	}

	log.Println(" Connected to PostgreSQL")
	return db
}
