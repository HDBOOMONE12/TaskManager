package handlers

import (
	"errors"
	"net/http"
	"strings"
	"time"

	"github.com/HDBOOMONE12/TaskManager/internal/service"
)

type TaskResponse struct {
	ID          int        `json:"id"`
	UserID      int        `json:"user_id"`
	Title       string     `json:"title"`
	Description string     `json:"description"`
	Status      string     `json:"status"`
	Priority    int        `json:"priority"`
	DueAt       *time.Time `json:"due_at"`
	CreatedAt   time.Time  `json:"created_at"`
	UpdatedAt   time.Time  `json:"updated_at"`
}

type CreateTaskRequest struct {
	Title       string `json:"title"`
	Description string `json:"description"`
	Status      string `json:"status"`
	Priority    int    `json:"priority"`
	DueAt       string `json:"due_at"`
}

func UserTasksHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		uid, perr := parseUserTasksListPath(r)
		if perr != nil {
			if errors.Is(perr, errBadPath) {
				http.NotFound(w, r)
				return
			}
			if errors.Is(perr, errBadID) {
				errorJSON(w, http.StatusBadRequest, "invalid user id")
				return
			}
			errorJSON(w, http.StatusBadRequest, "bad request")
			return
		}
		if _, ok := service.GetUserByID(uid); !ok {
			errorJSON(w, http.StatusNotFound, "user not found")
			return
		}

		list := service.ListTasksByUser(uid)
		resp := make([]TaskResponse, 0, len(list))
		for _, t := range list {
			resp = append(resp, TaskResponse{
				ID: t.ID, UserID: t.UserID, Title: t.Title, Description: t.Description,
				Status: t.Status, Priority: t.Priority, DueAt: t.DueAt,
				CreatedAt: t.CreatedAt, UpdatedAt: t.UpdatedAt,
			})
		}
		writeJSON(w, http.StatusOK, resp)

	case http.MethodPost:
		ct := r.Header.Get("Content-Type")
		if !strings.HasPrefix(ct, "application/json") {
			errorJSON(w, http.StatusUnsupportedMediaType, "Content-Type must be application/json")
			return
		}

		uid, perr := parseUserTasksListPath(r)
		if perr != nil {
			if errors.Is(perr, errBadPath) {
				http.NotFound(w, r)
				return
			}
			if errors.Is(perr, errBadID) {
				errorJSON(w, http.StatusBadRequest, "invalid user id")
				return
			}
			errorJSON(w, http.StatusBadRequest, "bad request")
			return
		}

		var req CreateTaskRequest
		if err := decodeJSON(w, r, &req, 1<<20); err != nil {
			respondDecodeError(w, err)
			return
		}

		var duePtr *time.Time
		if s := strings.TrimSpace(req.DueAt); s != "" {
			t, err := time.Parse(time.RFC3339, s)
			if err != nil {
				errorJSON(w, http.StatusBadRequest, "invalid due_at, use RFC3339 e.g 2025-08-20T10:00:00Z")
				return
			}
			duePtr = &t
		}

		task, err := service.CreateTask(uid, req.Title, req.Description, req.Status, req.Priority, duePtr)
		if err != nil {
			respondTaskError(w, err)
			return
		}

		writeJSON(w, http.StatusCreated, TaskResponse{
			ID: task.ID, UserID: task.UserID, Title: task.Title, Description: task.Description,
			Status: task.Status, Priority: task.Priority, DueAt: task.DueAt,
			CreatedAt: task.CreatedAt, UpdatedAt: task.UpdatedAt,
		})

	default:
		w.Header().Set("Allow", "GET, POST")
		errorJSON(w, http.StatusMethodNotAllowed, "method not allowed")
	}
}

func UserTaskDetailHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		uid, tid, perr := parseUserTaskDetailPath(r)
		if perr != nil {
			if errors.Is(perr, errBadPath) {
				http.NotFound(w, r)
				return
			}
			if errors.Is(perr, errBadID) {
				errorJSON(w, http.StatusBadRequest, "invalid user id")
				return
			}
			if errors.Is(perr, errBadTaskID) {
				errorJSON(w, http.StatusBadRequest, "invalid task id")
				return
			}
			errorJSON(w, http.StatusBadRequest, "bad request")
			return
		}
		if _, ok := service.GetUserByID(uid); !ok {
			errorJSON(w, http.StatusNotFound, "user not found")
			return
		}
		task, ok := service.GetTaskByUser(uid, tid)
		if !ok {
			errorJSON(w, http.StatusNotFound, "task not found")
			return
		}

		writeJSON(w, http.StatusOK, TaskResponse{
			ID: task.ID, UserID: task.UserID, Title: task.Title, Description: task.Description,
			Status: task.Status, Priority: task.Priority, DueAt: task.DueAt,
			CreatedAt: task.CreatedAt, UpdatedAt: task.UpdatedAt,
		})

	case http.MethodPut:
		ct := r.Header.Get("Content-Type")
		if !strings.HasPrefix(ct, "application/json") {
			errorJSON(w, http.StatusUnsupportedMediaType, "Content-Type must be application/json")
			return
		}

		var req struct {
			Title       string `json:"title"`
			Description string `json:"description"`
			Status      string `json:"status"`
			Priority    int    `json:"priority"`
			DueAt       string `json:"due_at"`
		}
		if err := decodeJSON(w, r, &req, 1<<20); err != nil {
			respondDecodeError(w, err)
			return
		}

		var duePtr *time.Time
		if s := strings.TrimSpace(req.DueAt); s != "" {
			t, err := time.Parse(time.RFC3339, s)
			if err != nil {
				errorJSON(w, http.StatusBadRequest, "invalid due_at, use RFC3339 e.g 2025-08-20T10:00:00Z")
				return
			}
			duePtr = &t
		}

		uid, tid, perr := parseUserTaskDetailPath(r)
		if perr != nil {
			if errors.Is(perr, errBadPath) {
				http.NotFound(w, r)
				return
			}
			if errors.Is(perr, errBadID) {
				errorJSON(w, http.StatusBadRequest, "invalid user id")
				return
			}
			if errors.Is(perr, errBadTaskID) {
				errorJSON(w, http.StatusBadRequest, "invalid task id")
				return
			}
			errorJSON(w, http.StatusBadRequest, "bad request")
			return
		}

		updated, err := service.UpdateTask(uid, tid, req.Title, req.Description, req.Status, req.Priority, duePtr)
		if err != nil {
			respondTaskError(w, err)
			return
		}

		writeJSON(w, http.StatusOK, TaskResponse{
			ID: updated.ID, UserID: updated.UserID, Title: updated.Title, Description: updated.Description,
			Status: updated.Status, Priority: updated.Priority, DueAt: updated.DueAt,
			CreatedAt: updated.CreatedAt, UpdatedAt: updated.UpdatedAt,
		})

	case http.MethodPatch:
		ct := r.Header.Get("Content-Type")
		if !strings.HasPrefix(ct, "application/json") {
			errorJSON(w, http.StatusUnsupportedMediaType, "Content-Type must be application/json")
			return
		}

		var req struct {
			Title       *string `json:"title"`
			Description *string `json:"description"`
			Status      *string `json:"status"`
			Priority    *int    `json:"priority"`
			DueAt       *string `json:"due_at"`
		}
		if err := decodeJSON(w, r, &req, 1<<20); err != nil {
			respondDecodeError(w, err)
			return
		}

		var dueAtProvided bool
		var duePtr *time.Time
		if req.DueAt == nil {
			dueAtProvided = false
		} else if strings.TrimSpace(*req.DueAt) == "" {
			dueAtProvided = true
			req.DueAt = nil
		} else {
			s := strings.TrimSpace(*req.DueAt)
			t, err := time.Parse(time.RFC3339, s)
			if err != nil {
				errorJSON(w, http.StatusBadRequest, "invalid due_at, use RFC3339 e.g 2025-08-20T10:00:00Z")
				return
			}
			dueAtProvided = true
			duePtr = &t
		}

		uid, tid, perr := parseUserTaskDetailPath(r)
		if perr != nil {
			if errors.Is(perr, errBadPath) {
				http.NotFound(w, r)
				return
			}
			if errors.Is(perr, errBadID) {
				errorJSON(w, http.StatusBadRequest, "invalid user id")
				return
			}
			if errors.Is(perr, errBadTaskID) {
				errorJSON(w, http.StatusBadRequest, "invalid task id")
				return
			}
			errorJSON(w, http.StatusBadRequest, "bad request")
			return
		}

		task, err := service.PatchTask(uid, tid, req.Title, req.Description, req.Status, req.Priority, dueAtProvided, duePtr)
		if err != nil {
			respondTaskError(w, err)
			return
		}

		writeJSON(w, http.StatusOK, TaskResponse{
			ID: task.ID, UserID: task.UserID, Title: task.Title, Description: task.Description,
			Status: task.Status, Priority: task.Priority, DueAt: task.DueAt,
			CreatedAt: task.CreatedAt, UpdatedAt: task.UpdatedAt,
		})

	case http.MethodDelete:
		uid, tid, perr := parseUserTaskDetailPath(r)
		if perr != nil {
			if errors.Is(perr, errBadPath) {
				http.NotFound(w, r)
				return
			}
			if errors.Is(perr, errBadID) {
				errorJSON(w, http.StatusBadRequest, "invalid user id")
				return
			}
			if errors.Is(perr, errBadTaskID) {
				errorJSON(w, http.StatusBadRequest, "invalid task id")
				return
			}
			errorJSON(w, http.StatusBadRequest, "bad request")
			return
		}

		ok := service.DeleteTaskByUser(uid, tid)
		if !ok {
			errorJSON(w, http.StatusNotFound, "task not found")
			return
		}

		w.WriteHeader(http.StatusNoContent)

	default:
		w.Header().Set("Allow", "GET, PUT, PATCH, DELETE")
		errorJSON(w, http.StatusMethodNotAllowed, "method not allowed")
	}
}
