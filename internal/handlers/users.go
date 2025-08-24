package handlers

import (
	"database/sql"
	"errors"
	"net/http"
	"strings"

	"github.com/HDBOOMONE12/TaskManager/internal/service"
)

type CreateUserRequest struct {
	Name  string `json:"name"`
	Email string `json:"email"`
}

type UserResponse struct {
	ID    int64  `json:"id"`
	Name  string `json:"name"`
	Email string `json:"email"`
}

var userSvc *service.UserService

func SetUserService(s *service.UserService) {
	userSvc = s
}

func UsersHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodHead:
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)

	case http.MethodGet:
		ctx := r.Context()
		list, err := userSvc.ListUsers(ctx)
		if err != nil {
			errorJSON(w, http.StatusInternalServerError, "internal error")
			return
		}
		resp := make([]UserResponse, 0, len(list))
		for _, u := range list {
			resp = append(resp, UserResponse{ID: u.ID, Name: u.Username, Email: u.Email})
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

		ctx := r.Context()
		u, err := userSvc.CreateUser(ctx, req.Name, req.Email)
		if err != nil {
			switch {
			case errors.Is(err, service.ErrEmptyName), errors.Is(err, service.ErrEmptyEmail):
				errorJSON(w, http.StatusBadRequest, err.Error())
			default:
				errorJSON(w, http.StatusInternalServerError, "internal error")
			}
			return
		}

		writeJSON(w, http.StatusCreated, UserResponse{ID: u.ID, Name: u.Username, Email: u.Email})

	default:
		w.Header().Set("Allow", "HEAD, GET, POST")
		errorJSON(w, http.StatusMethodNotAllowed, "method not allowed")
	}
}

func UserDetailHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodHead:
		w.Header().Set("Content-Type", "application/json")
		id, perr := parseUserID(r)
		if perr != nil {
			if errors.Is(perr, errBadPath) {
				w.WriteHeader(http.StatusNotFound)
				return
			}
			if errors.Is(perr, errBadID) {
				w.WriteHeader(http.StatusBadRequest)
				return
			}
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		ctx := r.Context()
		_, err := userSvc.GetUserByID(ctx, int64(id))
		if err != nil {
			w.WriteHeader(http.StatusNotFound)
			return
		}
		w.WriteHeader(http.StatusOK)

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

		ctx := r.Context()
		u, err := userSvc.GetUserByID(ctx, int64(id))
		if err != nil {
			// считаем любую ошибку — как not found, либо можешь расширить проверку:
			// if errors.Is(err, sql.ErrNoRows) { ... }
			errorJSON(w, http.StatusNotFound, "user not found")
			return
		}
		writeJSON(w, http.StatusOK, UserResponse{ID: u.ID, Name: u.Username, Email: u.Email})
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

		ctx := r.Context()
		u, err := userSvc.UpdateUserByID(ctx, int64(id), req.Name, req.Email)
		if err != nil {
			switch {
			case errors.Is(err, service.ErrUserNotFound), errors.Is(err, sql.ErrNoRows):
				errorJSON(w, http.StatusNotFound, "user not found")
			case errors.Is(err, service.ErrEmptyName), errors.Is(err, service.ErrEmptyEmail):
				errorJSON(w, http.StatusBadRequest, err.Error())
			default:
				errorJSON(w, http.StatusInternalServerError, "internal error")
			}
			return
		}

		writeJSON(w, http.StatusOK, UserResponse{ID: u.ID, Name: u.Username, Email: u.Email})
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

		ctx := r.Context()
		u, err := userSvc.PatchUserByID(ctx, int64(id), req.Name, req.Email)
		if err != nil {
			switch {
			case errors.Is(err, service.ErrUserNotFound), errors.Is(err, sql.ErrNoRows):
				errorJSON(w, http.StatusNotFound, "user not found")
			case errors.Is(err, service.ErrEmptyName), errors.Is(err, service.ErrEmptyEmail):
				errorJSON(w, http.StatusBadRequest, err.Error())
			default:
				errorJSON(w, http.StatusInternalServerError, "internal error")
			}
			return
		}

		writeJSON(w, http.StatusOK, UserResponse{ID: u.ID, Name: u.Username, Email: u.Email})
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

		ctx := r.Context()
		if err := userSvc.DeleteUserByID(ctx, int64(id)); err != nil {
			switch {
			case errors.Is(err, service.ErrUserNotFound), errors.Is(err, sql.ErrNoRows):
				errorJSON(w, http.StatusNotFound, "user not found")
			default:
				errorJSON(w, http.StatusInternalServerError, "internal error")
			}
			return
		}
		w.WriteHeader(http.StatusNoContent)
		return

	default:
		w.Header().Set("Allow", "HEAD, GET, PUT, PATCH, DELETE")
		errorJSON(w, http.StatusMethodNotAllowed, "method not allowed")
	}
}
