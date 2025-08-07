package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strings"
)

var nextId int = 1

func main() {
	http.HandleFunc("/users", createUserHandler)
	log.Fatal(http.ListenAndServe(":8080", nil))

}

type User struct {
	ID    int    `json:"-"`
	Name  string `json:"name"`
	Email string `json:"email"`
}

type userResponse struct {
	ID    int    `json:"id"`
	Name  string `json:"name"`
	Email string `json:"email"`
}

func createUserHandler(w http.ResponseWriter, r *http.Request) {
	ct := r.Header.Get("Content-Type")
	if !strings.HasPrefix(ct, "application/json") {
		w.WriteHeader(http.StatusUnsupportedMediaType)
		fmt.Fprintln(w, "Content-Type must be application/json")
		return
	}

	r.Body = http.MaxBytesReader(w, r.Body, 1<<20)

	var user User
	if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
		if strings.Contains(err.Error(), "request body too large") {
			w.WriteHeader(http.StatusRequestEntityTooLarge)
			fmt.Fprintln(w, "request body too largee")
			return
		}

		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintln(w, "Invalid JSON")
		return
	}

	if user.Name == "" || user.Email == "" {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintln(w, "missing name or email")
		return
	}
	user.ID = nextId
	nextId++

	resp := userResponse{
		ID:    user.ID,
		Name:  user.Name,
		Email: user.Email,
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(resp)

}
