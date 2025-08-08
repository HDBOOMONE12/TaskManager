package handlers

import (
	"encoding/json"
	"errors"
	"github.com/HDBOOMONE12/TaskManager/internal/service"
	"net/http"
	"strings"
)

type CreateUserRequest struct {
	Name  string `json:"name"`
	Email string `json:"email"`
}

type UserResponse struct {
	ID    int    `json:"id"`
	Name  string `json:"name"`
	Email string `json:"email"`
}

func UsersHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		list := service.ListUsers()
		respUsers := make([]UserResponse, 0, len(list))
		for _, user := range list {
			respUsers = append(respUsers, UserResponse{
				ID: user.ID, Name: user.Name, Email: user.Email,
			})
		}
		writeJSON(w, http.StatusOK, respUsers)

	case http.MethodPost:
		contentType := r.Header.Get("Content-Type")
		if !strings.HasPrefix(contentType, "application/json") {
			errorJSON(w, http.StatusUnsupportedMediaType, "Content-Type must be application/json")
			return
		}

		r.Body = http.MaxBytesReader(w, r.Body, 1<<20)

		var req CreateUserRequest
		dec := json.NewDecoder(r.Body)
		dec.DisallowUnknownFields()
		if err := dec.Decode(&req); err != nil {
			if strings.Contains(err.Error(), "request body too large") {
				errorJSON(w, http.StatusRequestEntityTooLarge, "request body too large")
				return
			}
			errorJSON(w, http.StatusBadRequest, "invalid JSON")
			return
		}

		if req.Name == "" || req.Email == "" {
			errorJSON(w, http.StatusBadRequest, "missing name or email")
			return
		}

		u := service.CreateUser(req.Name, req.Email)
		writeJSON(w, http.StatusCreated, UserResponse{
			ID: u.ID, Name: u.Name, Email: u.Email,
		})
	default:
		w.Header().Set("Allow", "GET, POST")
		errorJSON(w, http.StatusMethodNotAllowed, "method not allowed")

	}
}

func UserDetailHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		id, perr := parseUserID(r)
		if perr != nil {
			if errors.Is(perr, errBadPath) {
				http.NotFound(w, r)
				return
			}
			if errors.Is(perr, errBadID) {
				errorJSON(w, http.StatusBadRequest, "invalid id")
				return
			}
			errorJSON(w, http.StatusBadRequest, "bad request")
			return
		}

		u, ok := service.GetUserByID(id)
		if !ok {
			errorJSON(w, http.StatusNotFound, "user not found")
			return
		}

		writeJSON(w, http.StatusOK, UserResponse{
			ID: u.ID, Name: u.Name, Email: u.Email,
		})
		return

	default:
		w.Header().Set("Allow", "GET")
		errorJSON(w, http.StatusMethodNotAllowed, "method not allowed")
		return
	}

}
