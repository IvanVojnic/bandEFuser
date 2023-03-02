package rpc

import (
	"context"
	"fmt"

	"github.com/IvanVojnic/bandEFuser/models"
	pr "github.com/IvanVojnic/bandEFuser/proto"

	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
)

type Auth interface {
	SignUp(ctx context.Context, user *models.User) error
	UpdateRefreshToken(ctx context.Context, rt string, id uuid.UUID) error
	SignIn(ctx context.Context, user *models.User) (models.Tokens, error)
}

type UserAuthServer struct {
	pr.UnimplementedUserAuthServer
	authServ Auth
}

func NewUserAuthServer(authServ Auth) *UserAuthServer {
	return &UserAuthServer{authServ: authServ}
}

func (s *UserAuthServer) SignUp(ctx context.Context, req *pr.SignUpRequest) (*pr.SignUpResponse, error) {
	user := models.User{ID: uuid.New(), Email: req.GetEmail(), Name: req.GetName(), Password: req.GetPassword()}
	err := s.authServ.SignUp(ctx, &user)
	if err != nil {
		logrus.WithFields(logrus.Fields{
			"Error creating user": err,
			"user":                user,
		}).Errorf("error while creating, %s", err)
		return &pr.SignUpResponse{IsCreated: false}, fmt.Errorf("error while creating user")
	}
	return &pr.SignUpResponse{IsCreated: true}, nil
}

func (s *UserAuthServer) SignIn(ctx context.Context, req *pr.SignInRequest) (*pr.SignInResponse, error) {
	password := req.GetPassword()
	user := models.User{Name: req.GetName(), Password: password}
	tokens, err := s.authServ.SignIn(ctx, &user)
	if err != nil {
		logrus.WithFields(logrus.Fields{
			"Error logging user": err,
			"user":               user,
		}).Errorf("error while logging, %s", err)
		return &pr.SignInResponse{At: "", Rt: ""}, fmt.Errorf("error while login user, %s", err)
	}
	return &pr.SignInResponse{At: tokens.AccessToken, Rt: tokens.RefreshToken}, fmt.Errorf("error while loogin, wrong credentials %s", err)
}
