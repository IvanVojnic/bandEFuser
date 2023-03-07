// Package service define auth services methods
package service

import (
	"context"
	"fmt"
	"time"

	"github.com/IvanVojnic/bandEFuser/models"

	"github.com/golang-jwt/jwt"
	"github.com/google/uuid"

	"golang.org/x/crypto/bcrypt"
)

// Auth interface consists of methods to auth user
type Auth interface {
	SignUp(ctx context.Context, user *models.User) error
	UpdateRefreshToken(ctx context.Context, rt string, id uuid.UUID) error
	SignIn(ctx context.Context, user models.Login) (*models.User, error)
}

// AuthServer define service user auth struct
type AuthServer struct {
	authRepo Auth
}

// NewUserAuthServer used to init auth struct
func NewUserAuthServer(authRepo Auth) *AuthServer {
	return &AuthServer{authRepo: authRepo}
}

// SignUp used to create user
func (s *AuthServer) SignUp(ctx context.Context, user *models.User) error {
	err := s.authRepo.SignUp(ctx, user)
	if err != nil {
		return fmt.Errorf("error while sign up, %s", err)
	}
	return nil
}

// SignIn used to sign in user
func (s *AuthServer) SignIn(ctx context.Context, login *models.Login) (models.Tokens, error) {
	userDB, err := s.authRepo.SignIn(ctx, *login)
	if err != nil {
		return models.Tokens{}, fmt.Errorf("error while login user, %s", err)
	}
	errPasswordCheck := CheckPasswordHash(login.Password, userDB.Password)
	if errPasswordCheck != nil {
		return models.Tokens{}, fmt.Errorf("error while loogin, wrong credentials %s", err)
	}
	rt, errRT := generateToken(userDB.ID, TokenRTDuration, SigningKeyRefresh)
	if errRT != nil {
		return models.Tokens{}, fmt.Errorf("error while generate rt, %s", err)
	}
	at, errAT := generateToken(userDB.ID, TokenATDuration, SigningKeyAccess)
	if errAT != nil {
		return models.Tokens{}, fmt.Errorf("error while generate at, %s", err)
	}
	errUpdateRT := s.authRepo.UpdateRefreshToken(ctx, rt, userDB.ID)
	if errUpdateRT != nil {
		return models.Tokens{}, fmt.Errorf("error while set rt, %s", err)
	}
	return models.Tokens{RefreshToken: rt, AccessToken: at}, nil
}

// UpdateRefreshToken used to update token
func (s *AuthServer) UpdateRefreshToken(ctx context.Context, rt string, id uuid.UUID) error {
	return nil
}

// CheckPasswordHash used to compare hashed and not hashed passwords
func CheckPasswordHash(password, hash string) error {
	return bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
}

// SigningKeyAccess is a secret key for tokens
const SigningKeyAccess = "al5jkvkls83l9cw6l"

// SigningKeyRefresh is a secret key for tokens
const SigningKeyRefresh = "jkvf7834lkjbas98"

// TokenRTDuration is a duration of rt life
const TokenRTDuration = 1 * time.Hour

// TokenATDuration is a duration of at life
const TokenATDuration = 100 * time.Minute

// generateToken used to generate tokens with id
func generateToken(id uuid.UUID, tokenDuration time.Duration, key string) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, &models.TokenClaims{
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Add(tokenDuration).Unix(),
			IssuedAt:  time.Now().Unix(),
		},
		UserID: id,
	})
	return token.SignedString([]byte(key))
}
