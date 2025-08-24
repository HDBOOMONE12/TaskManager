package service

import (
	"context"
	"errors"
	"github.com/HDBOOMONE12/TaskManager/internal/entity"
	"github.com/HDBOOMONE12/TaskManager/internal/storage"
)

var (
	ErrEmptyName  = errors.New("empty name")
	ErrEmptyEmail = errors.New("empty email")
)

type UserService struct {
	repo *storage.UserRepo
}

func NewUserService(repo *storage.UserRepo) *UserService {
	return &UserService{repo: repo}
}

func (s *UserService) CreateUser(ctx context.Context, name, email string) (entity.User, error) {
	if name == "" {
		return entity.User{}, ErrEmptyName
	}
	if email == "" {
		return entity.User{}, ErrEmptyEmail
	}

	user := &entity.User{
		Username: name,
		Email:    email,
	}
	if err := s.repo.Create(ctx, user); err != nil {
		return entity.User{}, err
	}
	return *user, nil
}

func (s *UserService) ListUsers(ctx context.Context) ([]entity.User, error) {
	return s.repo.GetAll(ctx)
}

func (s *UserService) GetUserByID(ctx context.Context, id int64) (entity.User, error) {
	return s.repo.GetByID(ctx, id)
}

func (s *UserService) UpdateUserByID(ctx context.Context, id int64, name, email string) (entity.User, error) {
	if name == "" {
		return entity.User{}, ErrEmptyName
	}
	if email == "" {
		return entity.User{}, ErrEmptyEmail
	}
	return s.repo.Update(ctx, id, name, email)
}

func (s *UserService) PatchUserByID(ctx context.Context, id int64, name, email *string) (entity.User, error) {
	return s.repo.Patch(ctx, id, name, email)
}

func (s *UserService) DeleteUserByID(ctx context.Context, id int64) error {
	return s.repo.Delete(ctx, id)
}
