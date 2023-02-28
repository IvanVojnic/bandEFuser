package service

import (
	"context"
	"github.com/IvanVojnic/bandEFuser/models"
	"github.com/google/uuid"
)

type UserComm interface {
	GetFriends(ctx context.Context, userID uuid.UUID) ([]*models.User, error)
	SendFriendsRequest(ctx context.Context, userSender uuid.UUID, userReceiver uuid.UUID) error
	AcceptFriendsRequest(ctx context.Context, userSenderID uuid.UUID, userID uuid.UUID) error
	DeclineFriendsRequest(ctx context.Context, userSenderID uuid.UUID, userID uuid.UUID) error
	FindUser(ctx context.Context, userEmail string) (*models.User, error)
	GetRequest(ctx context.Context, userID uuid.UUID) ([]*models.User, error)
}

type UserCommServer struct {
	userCommRepo UserComm
}

func NewUserCommServer(userCommRepo UserComm) *UserCommServer {
	return &UserCommServer{userCommRepo: userCommRepo}
}

func (s *UserCommServer) GetFriends(ctx context.Context, userID uuid.UUID) ([]*models.User, error) {
	return s.userCommRepo.GetFriends(ctx, userID)
}

func (s *UserCommServer) SendFriendsRequest(ctx context.Context, userSenderID uuid.UUID, userReceiverID uuid.UUID) error {
	return s.userCommRepo.SendFriendsRequest(ctx, userSenderID, userReceiverID)
}

func (s *UserCommServer) AcceptFriendsRequest(ctx context.Context, userSenderID uuid.UUID, userID uuid.UUID) error {
	return s.userCommRepo.AcceptFriendsRequest(ctx, userSenderID, userID)
}

func (s *UserCommServer) DeclineFriendsRequest(ctx context.Context, userSenderID uuid.UUID, userID uuid.UUID) error {
	return s.userCommRepo.DeclineFriendsRequest(ctx, userSenderID, userID)
}

func (s *UserCommServer) FindUser(ctx context.Context, userEmail string) (*models.User, error) {
	return s.userCommRepo.FindUser(ctx, userEmail)
}

func (s *UserCommServer) GetRequest(ctx context.Context, userID uuid.UUID) ([]*models.User, error) {
	return s.userCommRepo.GetRequest(ctx, userID)
}
