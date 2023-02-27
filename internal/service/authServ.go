package service

import (
	"context"
	"fmt"
	"github.com/IvanVojnic/bandEFuser/internal/utils"
	"github.com/IvanVojnic/bandEFuser/models"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

type Auth interface {
	SignUp(ctx context.Context, user *models.User) error
	UpdateRefreshToken(ctx context.Context, rt string, id uuid.UUID) error
	GetUserByID(ctx context.Context, userID uuid.UUID) (models.User, error)
	SignIn(ctx context.Context, user *models.User) error
}

type AuthServer struct {
	authRepo Auth
}

func NewUserAuthServer(authRepo Auth) *AuthServer {
	return &AuthServer{authRepo: authRepo}
}

func (s *AuthServer) SignUp(ctx context.Context, user *models.User) error {
	err := s.authRepo.SignUp(ctx, user)
	if err != nil {
		return fmt.Errorf("error while sign up, %s", err)
	}
	return nil
}

func (s *AuthServer) SignIn(ctx context.Context, user *models.User) (models.Tokens, error) {
	password := user.Password
	err := s.authRepo.SignIn(ctx, user)
	if err != nil {
		return models.Tokens{}, fmt.Errorf("error while login user, %s", err)
	}
	match := CheckPasswordHash(password, user.Password)
	if match {
		rt, errRT := utils.GenerateToken(user.ID, utils.TokenRTDuration)
		if errRT != nil {
			return models.Tokens{}, fmt.Errorf("error while generate rt, %s", err)
		}
		at, errAT := utils.GenerateToken(user.ID, utils.TokenATDuretion)
		if errAT != nil {
			return models.Tokens{}, fmt.Errorf("error while generate at, %s", err)
		}
		errUpdateRT := s.authRepo.UpdateRefreshToken(ctx, rt, user.ID)
		if errUpdateRT != nil {
			return models.Tokens{}, fmt.Errorf("error while set rt, %s", err)
		}
		return models.Tokens{RefreshToken: rt, AccessToken: at}, nil
	}
	return models.Tokens{}, fmt.Errorf("error while loogin, wrong credentials %s", err)
}

func (s *AuthServer) UpdateRefreshToken(ctx context.Context, rt string, id uuid.UUID) error {
	return nil
}

func CheckPasswordHash(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}
