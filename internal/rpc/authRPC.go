package rpc

import (
	"context"
	"fmt"
	"github.com/IvanVojnic/bandEFuser/models"
	pr "github.com/IvanVojnic/bandEFuser/proto"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

type Auth interface {
	SignUp(ctx context.Context, user *models.User) error
	UpdateRefreshToken(ctx context.Context, rt string, id uuid.UUID) error
	SignIn(ctx context.Context, user *models.User) (models.Tokens, error)
}

type UserComm interface {
	GetFriends(ctx context.Context, userID uuid.UUID, users []*pr.User) error
	SendFriendsRequest(ctx context.Context, userSender uuid.UUID, userReceiver uuid.UUID) error
	AcceptFriendsRequest(ctx context.Context, userSenderID uuid.UUID, userID uuid.UUID) error
	DeclineFriendsRequest(ctx context.Context, userSenderID uuid.UUID, userID uuid.UUID) error
	FindUser(ctx context.Context, userEmail string) (*pr.User, error)
	GetRequest(ctx context.Context, userID uuid.UUID, users []*pr.User) error
}

type UserServer struct {
	pr.UnimplementedUserServer
	authServ     Auth
	userCommServ UserComm
}

func NewUserAuthServer(authServ Auth, userCommServ UserComm) *UserServer {
	return &UserServer{authServ: authServ, userCommServ: userCommServ}
}

func (s *UserServer) SignUp(ctx context.Context, req *pr.SignUpRequest) (*pr.SignUpResponse, error) {
	user := models.User{ID: uuid.New(), Email: req.GetEmail(), Name: req.GetName(), Password: req.GetPassword()}
	err := s.authServ.SignUp(ctx, &user)
	if err != nil {
		return &pr.SignUpResponse{IsCreated: false}, fmt.Errorf("error while creating user")
	}
	return &pr.SignUpResponse{IsCreated: true}, nil
}

func (s *UserServer) SignIn(ctx context.Context, req *pr.SignInRequest) (*pr.SignInResponse, error) {
	password := req.GetPassword()
	user := models.User{Name: req.GetName(), Password: password}
	tokens, err := s.authServ.SignIn(ctx, &user)
	if err != nil {
		return &pr.SignInResponse{At: "", Rt: ""}, fmt.Errorf("error while login user, %s", err)
	}
	return &pr.SignInResponse{At: tokens.AccessToken, Rt: tokens.RefreshToken}, fmt.Errorf("error while loogin, wrong credentials %s", err)
}

func CheckPasswordHash(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}
