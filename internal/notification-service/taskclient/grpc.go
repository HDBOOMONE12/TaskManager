package taskclient

import (
	"context"
	userspb "github.com/HDBOOMONE12/TaskManager/internal/taskmanager/proto"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"log"
)

type TaskGRPCClient struct {
	conn   *grpc.ClientConn
	client userspb.UserServiceClient
}

func NewTaskGRPCClient(addr string) (*TaskGRPCClient, error) {
	grpcConn, err := grpc.Dial(addr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("fail to connect: %v", err)
	}

	client := userspb.NewUserServiceClient(grpcConn)

	return &TaskGRPCClient{
		conn:   grpcConn,
		client: client,
	}, nil
}

func (c *TaskGRPCClient) HasUserWithEmail(ctx context.Context, email string) (bool, error) {
	request := userspb.EmailRequest{Email: email}
	response, err := c.client.HasUserWithEmail(ctx, &request)
	if err != nil {
		return false, err
	}

	return response.GetExists(), nil
}
