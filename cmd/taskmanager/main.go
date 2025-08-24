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

	"github.com/HDBOOMONE12/TaskManager/internal/db"
	"github.com/HDBOOMONE12/TaskManager/internal/handlers"
	"github.com/HDBOOMONE12/TaskManager/internal/service"
	"github.com/HDBOOMONE12/TaskManager/internal/storage"

	_ "github.com/jackc/pgx/v5/stdlib"
)

func main() {

	database := db.Init()
	defer database.Close()

	userRepo := storage.NewUserRepo(database)
	taskRepo := storage.NewTaskRepo(database)

	userSvc := service.NewUserService(userRepo)
	taskSvc := service.NewTaskService(taskRepo)

	handlers.SetUserService(userSvc)
	handlers.SetTaskService(taskSvc)

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
	mux.HandleFunc("/users", handlers.UsersHandler)
	mux.HandleFunc("/users/", handlers.UsersSubtreeHandler)
	return mux
}
