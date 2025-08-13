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
	_, ok := GetUserByID(userID)

	if !ok {
		return Task{}, ErrUserNotFound
	}

	if title == "" {
		return Task{}, ErrEmptyTitle
	}

	ok1 := isValidStatus(status)
	if !ok1 {
		return Task{}, ErrBadStatus
	}

	ok2 := isValidPriority(priority)
	if !ok2 {
		return Task{}, ErrBadPriority
	}

	tasksMu.Lock()
	defer tasksMu.Unlock()
	now := time.Now()
	createdAt := now
	updatedAt := now

	id := nextTaskID
	nextTaskID++

	task := Task{
		ID:          id,
		UserID:      userID,
		Title:       title,
		Description: desc,
		Status:      status,
		Priority:    priority,
		DueAt:       dueAt,
		CreatedAt:   createdAt,
		UpdatedAt:   updatedAt,
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

func DeleteTask(id int) bool {
	tasksMu.Lock()
	defer tasksMu.Unlock()
	_, ok := tasks[id]
	if !ok {
		return false
	}
	delete(tasks, id)
	return true
}
