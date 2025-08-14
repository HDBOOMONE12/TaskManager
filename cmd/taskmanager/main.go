package main

import (
	"github.com/HDBOOMONE12/TaskManager/internal/handlers"
	"log"
	"net/http"
)

func main() {
	mux := http.NewServeMux()
	mux.HandleFunc("/users", handlers.UsersHandler)
	mux.HandleFunc("/users/", handlers.UsersSubtreeHandler)
	log.Fatal(http.ListenAndServe(":8080", mux))
}
