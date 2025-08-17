package main

import (
	"context"
	"errors"
	"github.com/HDBOOMONE12/TaskManager/internal/handlers"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func main() {

	mux := buildMux()
	srv := &http.Server{
		Handler:           mux,
		Addr:              ":8080",
		ReadHeaderTimeout: 5 * time.Second,
		WriteTimeout:      30 * time.Second,
		ReadTimeout:       15 * time.Second,
		IdleTimeout:       60 * time.Second,
	}

	log.Printf("слушаю %s", srv.Addr)
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

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*15)
	defer cancel()

	err := srv.Shutdown(ctx)
	switch {
	case err == nil:
		log.Printf("graceful done")
	case errors.Is(err, context.DeadlineExceeded):
		log.Printf("stopped by deadline")
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
