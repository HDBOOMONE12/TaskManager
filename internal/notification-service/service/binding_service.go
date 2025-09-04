package service

import (
	"context"
	"github.com/HDBOOMONE12/TaskManager/internal/notification-service/storage"
)

type BindingService struct {
	repo *storage.TelegramBindingRepo
}

func NewBindingService(repo *storage.TelegramBindingRepo) *BindingService {
	return &BindingService{repo: repo}
}

func (s *BindingService) BindEmailToChat(ctx context.Context, email string, chatID int64) error {
	return s.repo.SaveBinding(ctx, email, chatID)
}

func (s *BindingService) GetChatIDByEmail(ctx context.Context, email string) (int64, error) {
	return s.repo.GetChatID(ctx, email)
}
