package handlers

import (
	"errors"
	"github.com/HDBOOMONE12/TaskManager/internal/taskmanager/service"
	"net/http"
)

func respondTaskError(w http.ResponseWriter, err error) bool {
	switch {
	case errors.Is(err, service.ErrUserNotFound):
		errorJSON(w, http.StatusNotFound, "user not found")
	case errors.Is(err, service.ErrEmptyTitle):
		errorJSON(w, http.StatusBadRequest, "invalid title")
	case errors.Is(err, service.ErrBadStatus):
		errorJSON(w, http.StatusBadRequest, "invalid task status")
	case errors.Is(err, service.ErrBadPriority):
		errorJSON(w, http.StatusBadRequest, "invalid task priority")
	case errors.Is(err, service.ErrTaskNotFound):
		errorJSON(w, http.StatusNotFound, "task not found")
	default:
		errorJSON(w, http.StatusInternalServerError, "internal server error")
	}
	return true
}
