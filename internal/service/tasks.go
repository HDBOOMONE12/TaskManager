package service

import (
	"context"
	"errors"
	"time"

	"github.com/HDBOOMONE12/TaskManager/internal/entity"
	"github.com/HDBOOMONE12/TaskManager/internal/storage"
)

const (
	StatusTodo       = "todo"
	StatusInProgress = "doing"
	StatusDone       = "done"

	MinPriority = 1
	MaxPriority = 5
)

var (
	ErrUserNotFound = errors.New("user not found")
	ErrEmptyTitle   = errors.New("empty title")
	ErrBadStatus    = errors.New("bad status")
	ErrBadPriority  = errors.New("bad priority")
	ErrTaskNotFound = errors.New("task not found")
)

type TaskService struct {
	repo *storage.TaskRepo
}

func NewTaskService(repo *storage.TaskRepo) *TaskService {
	return &TaskService{repo: repo}
}

func isValidStatus(status string) bool {
	switch status {
	case StatusTodo, StatusInProgress, StatusDone:
		return true
	}
	return false
}

func isValidPriority(priority int64) bool {
	return priority >= MinPriority && priority <= MaxPriority
}

func (s *TaskService) CreateTask(ctx context.Context, userID int64, title, desc, status string, priority int64, dueAt *time.Time) (entity.Task, error) {
	if status == "" {
		status = StatusTodo
	}
	if priority == 0 {
		priority = 3
	}
	if title == "" {
		return entity.Task{}, ErrEmptyTitle
	}
	if !isValidStatus(status) {
		return entity.Task{}, ErrBadStatus
	}
	if !isValidPriority(priority) {
		return entity.Task{}, ErrBadPriority
	}

	t := &entity.Task{
		UserID:      userID,
		Title:       title,
		Description: desc,
		Status:      status,
		Priority:    priority,
		DueAt:       dueAt,
	}
	if err := s.repo.Create(ctx, t); err != nil {
		return entity.Task{}, err
	}
	return *t, nil
}

func (s *TaskService) GetTaskByID(ctx context.Context, id int64) (entity.Task, error) {
	return s.repo.GetByID(ctx, id)
}

func (s *TaskService) ListTasksByUser(ctx context.Context, userID int64) ([]entity.Task, error) {
	return s.repo.GetByUserID(ctx, userID)
}

func (s *TaskService) UpdateTask(ctx context.Context, uid, tid int64, title, desc, status string, priority int64, dueAt *time.Time) (entity.Task, error) {
	if title == "" {
		return entity.Task{}, ErrEmptyTitle
	}
	if !isValidStatus(status) {
		return entity.Task{}, ErrBadStatus
	}
	if !isValidPriority(priority) {
		return entity.Task{}, ErrBadPriority
	}

	t := &entity.Task{
		ID:          tid,
		UserID:      uid,
		Title:       title,
		Description: desc,
		Status:      status,
		Priority:    priority,
		DueAt:       dueAt,
	}
	return s.repo.Update(ctx, t)
}

func (s *TaskService) PatchTask(
	ctx context.Context,
	uid, tid int64,
	title, desc, status *string,
	priority *int,
	dueAtProvided bool,
	dueAt *time.Time,
) (entity.Task, error) {
	cur, err := s.repo.GetByID(ctx, tid)
	if err != nil {
		return entity.Task{}, ErrTaskNotFound
	}
	if cur.UserID != uid {
		return entity.Task{}, ErrTaskNotFound
	}
	if title != nil && *title == "" {
		return entity.Task{}, ErrEmptyTitle
	}
	if status != nil && !isValidStatus(*status) {
		return entity.Task{}, ErrBadStatus
	}
	if priority != nil && !isValidPriority(int64(*priority)) {
		return entity.Task{}, ErrBadPriority
	}
	return s.repo.Patch(ctx, uid, tid, title, desc, status, priority, dueAtProvided, dueAt)
}

func (s *TaskService) DeleteTaskByUser(ctx context.Context, uid, tid int64) error {
	cur, err := s.repo.GetByID(ctx, tid)
	if err != nil {
		return ErrTaskNotFound
	}
	if cur.UserID != uid {
		return ErrTaskNotFound
	}
	return s.repo.Delete(ctx, tid)
}
