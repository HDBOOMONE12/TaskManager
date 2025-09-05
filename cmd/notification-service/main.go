package main

import (
	"context"
	"errors"
	"github.com/HDBOOMONE12/TaskManager/internal/notification-service/handlers"
	"github.com/HDBOOMONE12/TaskManager/internal/notification-service/senders"
	"github.com/HDBOOMONE12/TaskManager/internal/notification-service/taskclient"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	_ "github.com/jackc/pgx/v5/stdlib"

	"github.com/HDBOOMONE12/TaskManager/internal/notification-service/db"
	"github.com/HDBOOMONE12/TaskManager/internal/notification-service/service"
	"github.com/HDBOOMONE12/TaskManager/internal/notification-service/storage"
)

func main() {
	cfg := LoadConfig()

	dbConn := db.Init(cfg.DatabaseURL)
	defer dbConn.Close()

	telegramSender := senders.NewTelegramSender(cfg.TelegramToken)
	repo := storage.NewTelegramBindingRepo(dbConn)

	taskClient, err := taskclient.NewTaskGRPCClient(cfg.BaseUrl)
	if err != nil {
		log.Fatalf("gRPC client error: %v", err)
	}

	bindingService := service.NewBindingService(repo, taskClient)
	h := handlers.NewWebhookHandler(telegramSender, bindingService)

	mux := buildMux(h)

	srv := &http.Server{
		Addr:              ":" + cfg.Port,
		Handler:           mux,
		ReadHeaderTimeout: 5 * time.Second,
		WriteTimeout:      30 * time.Second,
		ReadTimeout:       15 * time.Second,
		IdleTimeout:       60 * time.Second,
	}

	log.Printf("Notification Service HTTP listening on %s", srv.Addr)
	go func() {
		if err := srv.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			log.Printf("http server error: %v", err)
		}
	}()

	log.Println("Notification Service started")

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
	switch {
	case err == nil:
		log.Printf("graceful shutdown complete")
	case errors.Is(err, context.DeadlineExceeded):
		log.Printf("shutdown deadline exceeded")
	default:
		log.Printf("shutdown error: %v", err)
	}
}

func buildMux(h http.Handler) *http.ServeMux {
	mux := http.NewServeMux()
	mux.Handle("/webhook", h)
	return mux
}
