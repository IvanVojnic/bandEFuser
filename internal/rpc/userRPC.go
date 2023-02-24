package rpc

import (
	"context"
	"github.com/IvanVojnic/bandEFuser/models"
	pr "github.com/IvanVojnic/bandEFuser/proto"
	"github.com/google/uuid"
)

type UserComm interface {
	GetFriends(ctx context.Context, userID uuid.UUID) ([]models.User, error)
	SendFriendsRequest(ctx context.Context, userSender uuid.UUID, userReceiver uuid.UUID) error
	AcceptFriendsRequest(ctx context.Context, userSenderID uuid.UUID, userID uuid.UUID) error
	DeclineFriendsRequest(ctx context.Context, userSenderID uuid.UUID, userID uuid.UUID) error
	FindUser(ctx context.Context, userEmail string) (models.User, error)
	GetRequest(ctx context.Context, userID uuid.UUID) ([]models.User, error)
}

type UserCommServer struct {
	pr.UserServer
	userCommRepo UserComm
}

func NewUserCommServer(userCommRepo UserComm) *UserCommServer {
	return &UserCommServer{userCommRepo: userCommRepo}
}

/*
func (s *UserCommServer) GetFriends(ctx context.Context, req *pr.GetFriendsRequest) (*pr.GetFriendsResponse, error) {
}

func (s *UserCommServer) SignIn(ctx context.Context, req *pr.SignUpRequest) (*pr.SignInResponse, error) {
}*/
