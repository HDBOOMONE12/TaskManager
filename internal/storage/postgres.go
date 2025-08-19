package storage

import (
	"context"
	"database/sql"
	"fmt"
	"time"
)

func NewDB(dsn string) (*sql.DB, error) {
	db, err := sql.Open("pgx", dsn)
	if err != nil {
		return nil, fmt.Errorf("open db: %w", err)
	}
	db.SetMaxOpenConns(12)
	db.SetMaxIdleConns(7)
	db.SetConnMaxLifetime(6 * time.Minute)
	db.SetConnMaxIdleTime(10 * time.Minute)
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*2)
	defer cancel()
	if err := db.PingContext(ctx); err != nil {
		_ = db.Close()
		return nil, fmt.Errorf("db ping failed: %w", err)
	}
	return db, nil

}
