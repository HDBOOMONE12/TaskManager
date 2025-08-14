package handlers

import (
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
		resp := make([]UserResponse, 0, len(list))
		for _, u := range list {
			resp = append(resp, UserResponse{ID: u.ID, Name: u.Name, Email: u.Email})
		}
		writeJSON(w, http.StatusOK, resp)

	case http.MethodPost:
		ct := r.Header.Get("Content-Type")
		if !strings.HasPrefix(ct, "application/json") {
			errorJSON(w, http.StatusUnsupportedMediaType, "Content-Type must be application/json")
			return
		}

		var req CreateUserRequest
		if err := decodeJSON(w, r, &req, 1<<20); err != nil {
			respondDecodeError(w, err)
			return
		}

		u, err := service.CreateUser(req.Name, req.Email)
		if err != nil {
			switch {
			case errors.Is(err, service.ErrEmptyName), errors.Is(err, service.ErrEmptyEmail):
				errorJSON(w, http.StatusBadRequest, err.Error())
			default:
				errorJSON(w, http.StatusInternalServerError, "internal error")
			}
			return
		}

		writeJSON(w, http.StatusCreated, UserResponse{ID: u.ID, Name: u.Name, Email: u.Email})

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
		writeJSON(w, http.StatusOK, UserResponse{ID: u.ID, Name: u.Name, Email: u.Email})
		return

	case http.MethodPut:
		ct := r.Header.Get("Content-Type")
		if !strings.HasPrefix(ct, "application/json") {
			errorJSON(w, http.StatusUnsupportedMediaType, "Content-Type must be application/json")
			return
		}

		var req CreateUserRequest
		if err := decodeJSON(w, r, &req, 1<<20); err != nil {
			respondDecodeError(w, err)
			return
		}

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

		u, err := service.UpdateUserByID(id, req.Name, req.Email)
		if err != nil {
			switch {
			case errors.Is(err, service.ErrUserNotFound):
				errorJSON(w, http.StatusNotFound, "user not found")
			case errors.Is(err, service.ErrEmptyName), errors.Is(err, service.ErrEmptyEmail):
				errorJSON(w, http.StatusBadRequest, err.Error())
			default:
				errorJSON(w, http.StatusInternalServerError, "internal error")
			}
			return
		}

		writeJSON(w, http.StatusOK, UserResponse{ID: u.ID, Name: u.Name, Email: u.Email})
		return

	case http.MethodPatch:
		ct := r.Header.Get("Content-Type")
		if !strings.HasPrefix(ct, "application/json") {
			errorJSON(w, http.StatusUnsupportedMediaType, "Content-Type must be application/json")
			return
		}

		var req struct {
			Name  *string `json:"name"`
			Email *string `json:"email"`
		}
		if err := decodeJSON(w, r, &req, 1<<20); err != nil {
			respondDecodeError(w, err)
			return
		}

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

		if req.Name == nil && req.Email == nil {
			errorJSON(w, http.StatusBadRequest, "no fields to update")
			return
		}

		u, err := service.PatchUserByID(id, req.Name, req.Email)
		if err != nil {
			switch {
			case errors.Is(err, service.ErrUserNotFound):
				errorJSON(w, http.StatusNotFound, "user not found")
			case errors.Is(err, service.ErrEmptyName), errors.Is(err, service.ErrEmptyEmail):
				errorJSON(w, http.StatusBadRequest, err.Error())
			default:
				errorJSON(w, http.StatusInternalServerError, "internal error")
			}
			return
		}

		writeJSON(w, http.StatusOK, UserResponse{ID: u.ID, Name: u.Name, Email: u.Email})
		return

	case http.MethodDelete:
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

		ok := service.DeleteUserByID(id)
		if !ok {
			errorJSON(w, http.StatusNotFound, "user not found")
			return
		}
		w.WriteHeader(http.StatusNoContent)
		return

	default:
		w.Header().Set("Allow", "GET, PUT, PATCH, DELETE")
		errorJSON(w, http.StatusMethodNotAllowed, "method not allowed")
	}
}
