package service

import (
	"errors"
	"sync"
	"time"
)

const (
	StatusTodo       = "todo"
	StatusInProgress = "in_progress"
	StatusDone       = "done"

	MinPriority = 1
	MaxPriority = 5
)

var (
	tasksMu    sync.RWMutex
	tasks      = make(map[int]Task)
	nextTaskID = 1
)

var (
	ErrUserNotFound = errors.New("user not found")
	ErrEmptyTitle   = errors.New("empty title")
	ErrBadStatus    = errors.New("bad status")
	ErrBadPriority  = errors.New("bad priority")
	ErrTaskNotFound = errors.New("task not found")
)

type Task struct {
	ID          int
	UserID      int
	Title       string
	Description string
	Status      string
	Priority    int
	DueAt       *time.Time
	CreatedAt   time.Time
	UpdatedAt   time.Time
}

func isValidStatus(status string) bool {
	switch status {
	case StatusTodo, StatusInProgress, StatusDone:
		return true
	}
	return false
}

func isValidPriority(priority int) bool {
	if priority >= MinPriority && priority <= MaxPriority {
		return true
	}
	return false
}

func CreateTask(userID int, title, desc, status string, priority int, dueAt *time.Time) (Task, error) {

	if _, ok := GetUserByID(userID); !ok {
		return Task{}, ErrUserNotFound
	}

	if status == "" {
		status = StatusTodo
	}
	if priority == 0 {
		priority = 3
	}

	if title == "" {
		return Task{}, ErrEmptyTitle
	}
	if !isValidStatus(status) {
		return Task{}, ErrBadStatus
	}
	if !isValidPriority(priority) {
		return Task{}, ErrBadPriority
	}

	tasksMu.Lock()
	defer tasksMu.Unlock()

	id := nextTaskID
	nextTaskID++

	now := time.Now()
	task := Task{
		ID:          id,
		UserID:      userID,
		Title:       title,
		Description: desc,
		Status:      status,
		Priority:    priority,
		DueAt:       dueAt,
		CreatedAt:   now,
		UpdatedAt:   now,
	}
	tasks[id] = task
	return task, nil
}

func GetTaskByID(id int) (Task, bool) {
	tasksMu.RLock()
	defer tasksMu.RUnlock()
	task, ok := tasks[id]
	return task, ok
}

func ListTasksByUser(userID int) []Task {
	tasksMu.RLock()
	defer tasksMu.RUnlock()
	tasksByUser := make([]Task, 0)
	for _, task := range tasks {
		if task.UserID == userID {
			tasksByUser = append(tasksByUser, task)
		}
	}
	return tasksByUser
}

func GetTaskByUser(uid, tid int) (Task, bool) {
	task, ok := GetTaskByID(tid)
	if !ok {
		return Task{}, false
	}
	if task.UserID != uid {
		return Task{}, false
	}
	return task, true
}

func UpdateTask(uid, tid int, title, desc, status string, priority int, dueAt *time.Time) (Task, error) {

	if title == "" {
		return Task{}, ErrEmptyTitle
	}
	if !isValidStatus(status) {
		return Task{}, ErrBadStatus
	}
	if !isValidPriority(priority) {
		return Task{}, ErrBadPriority
	}

	tasksMu.Lock()
	defer tasksMu.Unlock()

	t, ok := tasks[tid]

	if !ok {
		return Task{}, ErrTaskNotFound
	}

	if t.UserID != uid {
		return Task{}, ErrTaskNotFound
	}

	t.Title = title
	t.Description = desc
	t.Status = status
	t.Priority = priority
	t.DueAt = dueAt
	t.UpdatedAt = time.Now()

	tasks[tid] = t
	return t, nil
}

func PatchTask(uid, tid int, title, desc, status *string, priority *int, dueAtProvided bool, dueAt *time.Time) (Task, error) {
	tasksMu.Lock()
	defer tasksMu.Unlock()

	t, ok := tasks[tid]
	if !ok || t.UserID != uid {
		return Task{}, ErrTaskNotFound
	}

	if title != nil {
		if *title == "" {
			return Task{}, ErrEmptyTitle
		}
		t.Title = *title
	}
	if desc != nil {
		t.Description = *desc
	}
	if status != nil {
		if !isValidStatus(*status) {
			return Task{}, ErrBadStatus
		}
		t.Status = *status
	}
	if priority != nil {
		if !isValidPriority(*priority) {
			return Task{}, ErrBadPriority
		}
		t.Priority = *priority
	}
	if dueAtProvided {
		t.DueAt = dueAt
	}

	t.UpdatedAt = time.Now()
	tasks[tid] = t
	return t, nil
}

func DeleteTaskByUser(uid, tid int) bool {
	tasksMu.Lock()
	defer tasksMu.Unlock()

	t, ok := tasks[tid]
	if !ok || t.UserID != uid {
		return false
	}

	delete(tasks, tid)
	return true
}
