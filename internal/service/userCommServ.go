// Package service define communicate services methods
package service

import (
	"context"
	"fmt"

	"github.com/IvanVojnic/bandEFuser/models"

	"github.com/google/uuid"
)

// UserComm interface consists of methods to communicate with user repo
type UserComm interface {
	GetFriends(ctx context.Context, userID uuid.UUID) ([]*models.User, error)
	SendFriendsRequest(ctx context.Context, userSender uuid.UUID, userReceiver uuid.UUID) error
	AcceptFriendsRequest(ctx context.Context, userSenderID uuid.UUID, userID uuid.UUID) error
	DeclineFriendsRequest(ctx context.Context, userSenderID uuid.UUID, userID uuid.UUID) error
	FindUser(ctx context.Context, userEmail string) (*models.User, error)
	GetRequest(ctx context.Context, userID uuid.UUID) ([]*models.User, error)
	GetUsers(ctx context.Context, usersID []*uuid.UUID) ([]*models.User, error)
	StorageFriendsRequest(ctx context.Context, userSender, userReceiver models.User) error
	GetUserByID(ctx context.Context, userID uuid.UUID) (*models.User, error)
}

// UserCommServer define service user communicate struct
type UserCommServer struct {
	userCommRepo UserComm
}

// NewUserCommServer used to init service user communicate struct
func NewUserCommServer(userCommRepo UserComm) *UserCommServer {
	return &UserCommServer{userCommRepo: userCommRepo}
}

// GetFriends used to get friends by repo
func (s *UserCommServer) GetFriends(ctx context.Context, userID uuid.UUID) ([]*models.User, error) { // nolint:dupl, gocritic
	return s.userCommRepo.GetFriends(ctx, userID)
}

// SendFriendsRequest used to send request by repo
func (s *UserCommServer) SendFriendsRequest(ctx context.Context, userSenderID, userReceiverID uuid.UUID) error {
	err := s.userCommRepo.SendFriendsRequest(ctx, userSenderID, userReceiverID)
	if err != nil {
		return fmt.Errorf("error while sending invite, %s", err)
	}
	userSender, err := s.userCommRepo.GetUserByID(ctx, userSenderID)
	if err != nil {
		return fmt.Errorf("error while getting user sender, %s", err)
	}
	userReceiver, err := s.userCommRepo.GetUserByID(ctx, userSenderID)
	if err != nil {
		return fmt.Errorf("error while getting user receiver, %s", err)
	}
	return s.userCommRepo.StorageFriendsRequest(ctx, *userSender, *userReceiver)
}

// AcceptFriendsRequest used to accept request by repo
func (s *UserCommServer) AcceptFriendsRequest(ctx context.Context, userSenderID, userID uuid.UUID) error {
	return s.userCommRepo.AcceptFriendsRequest(ctx, userSenderID, userID)
}

// DeclineFriendsRequest used to decline request
func (s *UserCommServer) DeclineFriendsRequest(ctx context.Context, userSenderID, userID uuid.UUID) error {
	return s.userCommRepo.DeclineFriendsRequest(ctx, userSenderID, userID)
}

// FindUser used to find user
func (s *UserCommServer) FindUser(ctx context.Context, userEmail string) (*models.User, error) {
	return s.userCommRepo.FindUser(ctx, userEmail)
}

// GetRequest used to getting request to be a friend by repo
func (s *UserCommServer) GetRequest(ctx context.Context, userID uuid.UUID) ([]*models.User, error) { // nolint:dupl, gocritic
	return s.userCommRepo.GetRequest(ctx, userID)
}

// GetUsers used to getting users by theirs id by repo
func (s *UserCommServer) GetUsers(ctx context.Context, usersID []*uuid.UUID) ([]*models.User, error) {
	return s.userCommRepo.GetUsers(ctx, usersID)
}

func (s *UserCommServer) GetUser(ctx context.Context, userID uuid.UUID) (*models.User, error) {
	return s.userCommRepo.GetUserByID(ctx, userID)
}
