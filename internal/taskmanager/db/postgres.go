package db

import (
	"database/sql"
	"log"
	"time"
)

func Init(dsn string) *sql.DB {

	db, err := sql.Open("pgx", dsn)
	if err != nil {
		log.Fatalf("cannot open DB: %v", err)
	}
	db.SetMaxOpenConns(10)
	db.SetMaxIdleConns(5)
	db.SetConnMaxLifetime(time.Minute * 30)

	if err := db.Ping(); err != nil {
		log.Fatalf("cannot connect to DB: %v", err)
	}

	log.Println(" Connected to PostgreSQL")
	return db
}
