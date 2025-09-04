package service

import (
	"context"
	"github.com/HDBOOMONE12/TaskManager/internal/notification-service/notifyerrors"
	"github.com/HDBOOMONE12/TaskManager/internal/notification-service/storage"
	"github.com/HDBOOMONE12/TaskManager/internal/notification-service/taskclient"
)

type BindingService struct {
	repo       *storage.TelegramBindingRepo
	taskClient *taskclient.TaskClient
}

func NewBindingService(repo *storage.TelegramBindingRepo, taskClient *taskclient.TaskClient) *BindingService {
	return &BindingService{repo: repo, taskClient: taskClient}
}

func (s *BindingService) BindEmailToChat(ctx context.Context, email string, chatID int64) error {
	check, err := s.taskClient.HasUserWithEmail(ctx, email)

	if err != nil {
		return err
	}

	if !check {
		return notifyerrors.ErrUserNotFound
	}

	return s.repo.SaveBinding(ctx, email, chatID)
}

func (s *BindingService) GetChatIDByEmail(ctx context.Context, email string) (int64, error) {
	return s.repo.GetChatID(ctx, email)
}
