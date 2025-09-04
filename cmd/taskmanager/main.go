package main

import (
	"context"
	"errors"
	"github.com/HDBOOMONE12/TaskManager/internal/taskmanager/db"
	handlers2 "github.com/HDBOOMONE12/TaskManager/internal/taskmanager/handlers"
	service2 "github.com/HDBOOMONE12/TaskManager/internal/taskmanager/service"
	storage2 "github.com/HDBOOMONE12/TaskManager/internal/taskmanager/storage"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	_ "github.com/jackc/pgx/v5/stdlib"
)

func main() {

	config := LoadConfig()

	database := db.Init(config.DatabaseURL)
	defer database.Close()

	userRepo := storage2.NewUserRepo(database)
	taskRepo := storage2.NewTaskRepo(database)

	userSvc := service2.NewUserService(userRepo)
	taskSvc := service2.NewTaskService(taskRepo)

	handlers2.SetUserService(userSvc)
	handlers2.SetTaskService(taskSvc)

	mux := buildMux()
	srv := &http.Server{
		Handler:           mux,
		Addr:              "localhost:8080",
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

	log.Println("Server started")

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

	err := srv.Shutdown(ctx)
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
	mux.HandleFunc("/users", handlers2.UsersHandler)
	mux.HandleFunc("/users/", handlers2.UsersSubtreeHandler)
	return mux
}
