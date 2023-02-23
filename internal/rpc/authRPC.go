package rpc

import (
	"context"
	"fmt"
	"github.com/google/uuid"
	"userMS/models"
	pr "userMS/proto"
)

type Auth interface {
	SignUp(ctx context.Context, user *models.User) error
	UpdateRefreshToken(ctx context.Context, rt string, id uuid.UUID) error
	GetUserByID(ctx context.Context, userID uuid.UUID) (models.User, error)
	SignIn(ctx context.Context, user *models.User) error
}

type UserServer struct {
	pr.UserServer
	authR Auth
}

func NewMs1Server() *UserServer {
	return &UserServer{}
}

func (s *UserServer) SignUp(ctx context.Context, req *pr.SignUpRequest) (*pr.SignUpResponse, error) {
	user := models.User{ID: uuid.New(), Email: req.GetEmail(), Name: req.GetName(), Password: req.GetPassword()}
	err := s.authR.SignUp(ctx, &user)
	if err != nil {
		return &pr.SignUpResponse{IsCreated: false}, fmt.Errorf("error while creating user")
	}
	return &pr.SignUpResponse{IsCreated: true}, nil
}
