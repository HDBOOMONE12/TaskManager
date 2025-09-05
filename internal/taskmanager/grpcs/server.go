package grpcs

import (
	"context"
	"database/sql"
	"errors"
	userspb "github.com/HDBOOMONE12/TaskManager/internal/taskmanager/proto"
	"github.com/HDBOOMONE12/TaskManager/internal/taskmanager/service"
)

type GrpcServer struct {
	UserService *service.UserService
	userspb.UnimplementedUserServiceServer
}

var _ userspb.UserServiceServer = &GrpcServer{}

func (s *GrpcServer) HasUserWithEmail(ctx context.Context, req *userspb.EmailRequest) (*userspb.UserExistsResponse, error) {
	_, err := s.UserService.GetByEmail(ctx, req.Email)
	if errors.Is(err, sql.ErrNoRows) {
		return &userspb.UserExistsResponse{
			Exists: false,
		}, nil
	}
	if err != nil {
		return nil, err
	}

	return &userspb.UserExistsResponse{
		Exists: true,
	}, nil
}
