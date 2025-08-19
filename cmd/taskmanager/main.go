package main

import (
	"context"
	"errors"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/HDBOOMONE12/TaskManager/internal/config"
	"github.com/HDBOOMONE12/TaskManager/internal/handlers"
	"github.com/HDBOOMONE12/TaskManager/internal/storage"
	_ "github.com/jackc/pgx/v5/stdlib"
)

func main() {
	cfg := config.LoadConfig()

	if !cfg.EnableHTTP {
		log.Printf("HTTP is disabled by config — exiting without starting the server")
		return
	}
	if cfg.DatabaseURL == "" {
		log.Printf("HTTP is enabled, but DATABASE_URL is empty — exiting")
		return
	}

	log.Printf("Connecting to DB: %s", config.MaskDSN(cfg.DatabaseURL))
	db, err := storage.NewDB(cfg.DatabaseURL)
	if err != nil {
		log.Printf("DB connection failed: %v", err)
		return
	}
	log.Printf("Database is ready")

	mux := buildMux()
	srv := &http.Server{
		Handler:           mux,
		Addr:              cfg.HTTPAddr,
		ReadHeaderTimeout: 5 * time.Second,
		WriteTimeout:      30 * time.Second,
		ReadTimeout:       15 * time.Second,
		IdleTimeout:       60 * time.Second,
	}

	log.Printf("HTTP listening on %s", srv.Addr)
	go func() {
		if err := srv.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			log.Printf("http server error: %v", err)
		}
	}()

	ch := make(chan os.Signal, 1)
	signal.Notify(ch, os.Interrupt, syscall.SIGTERM)

	sig := <-ch
	log.Printf("signal received: %v — starting graceful shutdown", sig)

	go func() {
		<-ch
		log.Printf("second signal — forcing exit")
		os.Exit(1)
	}()

	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()

	err = srv.Shutdown(ctx)

	if closeErr := db.Close(); closeErr != nil {
		log.Printf("db close error: %v", closeErr)
	}

	switch {
	case err == nil:
		log.Printf("graceful shutdown complete")
	case errors.Is(err, context.DeadlineExceeded):
		log.Printf("shutdown deadline exceeded")
	default:
		log.Printf("shutdown error: %v", err)
	}
}

func buildMux() *http.ServeMux {
	mux := http.NewServeMux()
	mux.HandleFunc("/users", handlers.UsersHandler)
	mux.HandleFunc("/users/", handlers.UsersSubtreeHandler)
	return mux
}
