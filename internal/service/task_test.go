package service

import (
	"context"
	"errors"
	"github.com/HDBOOMONE12/TaskManager/internal/entity"
	"github.com/HDBOOMONE12/TaskManager/internal/mocks"
	"github.com/golang/mock/gomock"
	"testing"
)

func TestIsValidStatus(t *testing.T) {
	tests := []struct {
		name  string
		input string
		want  bool
	}{
		{"Valid todo status", StatusTodo, true},
		{"Valid doing status", StatusInProgress, true},
		{"Valid done status", StatusDone, true},
		{"Invalid Skebob status", "Skebob", false},
		{"Invalid empty status", "", false},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			got := isValidStatus(tt.input)
			if got != tt.want {
				t.Errorf("isValidStatus(%s): got %v, want %v", tt.input, got, tt.want)
			}
		})
	}
}

func TestIsValidPriority(t *testing.T) {
	tests := []struct {
		name  string
		input int64
		want  bool
	}{
		{"Valid priority 1", 1, true},
		{"Valid priority 5", 5, true},
		{"Invalid priority 6", 6, false},
		{"Invalid priority 0", 0, false},
		{"Valid priority 3", 3, true},
		{"Invalid negative priority", -1, false},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			got := isValidPriority(tt.input)
			if got != tt.want {
				t.Errorf("isValidPriority(%d): got %v, want %v", tt.input, got, tt.want)
			}
		})
	}
}

func TestTaskService_CreateTask(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mocks.NewMockTaskRepository(ctrl)
	svc := NewTaskService(mockRepo)

	dbErr := errors.New("DB error")

	tests := []struct {
		name      string
		title     string
		status    string
		priority  int64
		mockSetup func()
		wantErr   error
	}{
		{
			name:     "success",
			title:    "Task1",
			status:   StatusTodo,
			priority: 3,
			mockSetup: func() {
				mockRepo.EXPECT().
					Create(gomock.Any(), gomock.Any()).
					Return(nil)
			},
			wantErr: nil,
		},
		{
			name:      "empty title",
			title:     "",
			status:    StatusTodo,
			priority:  3,
			mockSetup: func() {},
			wantErr:   ErrEmptyTitle,
		},
		{
			name:      "invalid status",
			title:     "bob",
			status:    "skebob",
			priority:  3,
			mockSetup: func() {},
			wantErr:   ErrBadStatus,
		},
		{
			name:      "invalid priority",
			title:     "bob",
			status:    StatusInProgress,
			priority:  50,
			mockSetup: func() {},
			wantErr:   ErrBadPriority,
		},
		{
			name:     "repository error",
			title:    "Task 1",
			status:   StatusTodo,
			priority: 3,
			mockSetup: func() {
				mockRepo.EXPECT().
					Create(gomock.Any(), gomock.Any()).
					Return(dbErr)
			},
			wantErr: dbErr,
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			tt.mockSetup()

			_, err := svc.CreateTask(
				context.Background(),
				1, tt.title, "desc", tt.status, tt.priority, nil,
			)

			if (err == nil) != (tt.wantErr == nil) {
				t.Fatalf("expected error %v, got %v", tt.wantErr, err)
			}
			if err != nil && err.Error() != tt.wantErr.Error() {
				t.Errorf("expected error %v, got %v", tt.wantErr, err)
			}
		})
	}
}

func TestTaskService_UpdateTask(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mocks.NewMockTaskRepository(ctrl)
	svc := NewTaskService(mockRepo)

	dbErr := errors.New("DB error")

	tests := []struct {
		name      string
		uid       int64
		tid       int64
		title     string
		status    string
		priority  int64
		mockSetup func()
		wantErr   error
	}{
		{
			name:     "success",
			title:    "Task1",
			status:   StatusTodo,
			priority: 3,
			mockSetup: func() {
				mockRepo.EXPECT().
					Update(gomock.Any(), gomock.Any()).
					Return(entity.Task{
						ID:       1,
						UserID:   1,
						Title:    "1",
						Status:   StatusTodo,
						Priority: 5,
					}, nil).
					Times(1)
			},
			wantErr: nil,
		},

		{
			name:     "empty title",
			title:    "",
			status:   StatusTodo,
			priority: 3,
			mockSetup: func() {
				mockRepo.EXPECT().
					Update(gomock.Any(), gomock.Any()).
					Times(0)
			},

			wantErr: ErrEmptyTitle,
		},
		{
			name:     "invalid status",
			title:    "bob",
			status:   "skebob",
			priority: 3,
			mockSetup: func() {
				mockRepo.EXPECT().
					Update(gomock.Any(), gomock.Any()).
					Times(0)
			},
			wantErr: ErrBadStatus,
		},
		{
			name:     "invalid priority",
			title:    "bob",
			status:   StatusInProgress,
			priority: 50,
			mockSetup: func() {
				mockRepo.EXPECT().
					Update(gomock.Any(), gomock.Any()).
					Times(0)
			},
			wantErr: ErrBadPriority,
		},
		{
			name:     "repository error",
			title:    "Task 1",
			status:   StatusTodo,
			priority: 3,
			mockSetup: func() {
				mockRepo.EXPECT().
					Update(
						gomock.Any(),
						gomock.Any(),
					).Return(entity.Task{}, dbErr).Times(1)
			},
			wantErr: dbErr,
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			tt.mockSetup()
			_, err := svc.UpdateTask(context.Background(),
				tt.uid, tt.tid, tt.title, "", tt.status, tt.priority, nil)

			if (err == nil) != (tt.wantErr == nil) {
				t.Fatalf("expected error %v, got %v", tt.wantErr, err)
			}
			if err != nil && err.Error() != tt.wantErr.Error() {
				t.Errorf("expected error %v, got %v", tt.wantErr, err)
			}
		})
	}
}
