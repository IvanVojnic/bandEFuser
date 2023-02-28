package service

import (
	"context"
	"fmt"
	"github.com/IvanVojnic/bandEFuser/models"
	"github.com/google/uuid"
)

type UserComm interface {
	GetFriends(ctx context.Context, userID uuid.UUID, users []*models.User) error
	SendFriendsRequest(ctx context.Context, userSender uuid.UUID, userReceiver uuid.UUID) error
	AcceptFriendsRequest(ctx context.Context, userSenderID uuid.UUID, userID uuid.UUID) error
	DeclineFriendsRequest(ctx context.Context, userSenderID uuid.UUID, userID uuid.UUID) error
	FindUser(ctx context.Context, userEmail string) (*models.User, error)
	GetRequest(ctx context.Context, userID uuid.UUID, users []*models.User) error
}

type UserCommServer struct {
	userCommRepo UserComm
}

func NewUserCommServer(userCommRepo UserComm) *UserCommServer {
	return &UserCommServer{userCommRepo: userCommRepo}
}

func (s *UserCommServer) GetFriends(ctx context.Context, userID uuid.UUID) ([]*models.User, error) {
	var users []*models.User
	err := s.userCommRepo.GetFriends(ctx, userID, users)
	if err != nil {
		return users, fmt.Errorf("error while gettingg friends from db, %s", err)
	}
	return users, nil
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
	var users []*models.User
	err := s.userCommRepo.GetRequest(ctx, userID, users)
	if err != nil {
		return users, fmt.Errorf("error while getting requests to be a friend from db, %s", err)
	}
	return users, nil
}
