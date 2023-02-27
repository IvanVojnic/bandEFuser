package rpc

import (
	"context"
	"fmt"
	"github.com/IvanVojnic/bandEFuser/internal/utils"
	"github.com/IvanVojnic/bandEFuser/models"
	pr "github.com/IvanVojnic/bandEFuser/proto"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

type Auth interface {
	SignUp(ctx context.Context, user *models.User) error
	UpdateRefreshToken(ctx context.Context, rt string, id uuid.UUID) error
	GetUserByID(ctx context.Context, userID uuid.UUID) (models.User, error)
	SignIn(ctx context.Context, user *models.User) error
}

type UserComm interface {
	GetFriends(ctx context.Context, userID uuid.UUID, users []*pr.User) error
	SendFriendsRequest(ctx context.Context, userSender uuid.UUID, userReceiver uuid.UUID) error
	AcceptFriendsRequest(ctx context.Context, userSenderID uuid.UUID, userID uuid.UUID) error
	DeclineFriendsRequest(ctx context.Context, userSenderID uuid.UUID, userID uuid.UUID) error
	FindUser(ctx context.Context, userEmail string) (*pr.User, error)
	GetRequest(ctx context.Context, userID uuid.UUID, users []*pr.User) error
}

// Tokens used to define at and rt
type Tokens struct {
	AccessToken  string `json:"access"`
	RefreshToken string `json:"refresh"`
}

type UserServer struct {
	pr.UnimplementedUserServer
	authRepo     Auth
	userCommRepo UserComm
}

func NewUserAuthServer(authRepo Auth, userCommRepo UserComm) *UserServer {
	return &UserServer{authRepo: authRepo, userCommRepo: userCommRepo}
}

func (s *UserServer) SignUp(ctx context.Context, req *pr.SignUpRequest) (*pr.SignUpResponse, error) {
	user := models.User{ID: uuid.New(), Email: req.GetEmail(), Name: req.GetName(), Password: req.GetPassword()}
	err := s.authRepo.SignUp(ctx, &user)
	if err != nil {
		return &pr.SignUpResponse{IsCreated: false}, fmt.Errorf("error while creating user")
	}
	return &pr.SignUpResponse{IsCreated: true}, nil
}

func (s *UserServer) SignIn(ctx context.Context, req *pr.SignInRequest) (*pr.SignInResponse, error) {
	password := req.GetPassword()
	user := models.User{Name: req.GetName(), Password: password}
	err := s.authRepo.SignIn(ctx, &user)
	if err != nil {
		return &pr.SignInResponse{At: "", Rt: ""}, fmt.Errorf("error while login user, %s", err)
	}
	match := CheckPasswordHash(password, user.Password)
	if match {
		rt, errRT := utils.GenerateToken(user.ID, utils.TokenRTDuration)
		if errRT != nil {
			return &pr.SignInResponse{At: "", Rt: ""}, fmt.Errorf("error while generate rt, %s", err)
		}
		at, errAT := utils.GenerateToken(user.ID, utils.TokenATDuretion)
		if errAT != nil {
			return &pr.SignInResponse{At: "", Rt: ""}, fmt.Errorf("error while generate at, %s", err)
		}
		errUpdateRT := s.authRepo.UpdateRefreshToken(ctx, rt, user.ID)
		if errUpdateRT != nil {
			return &pr.SignInResponse{At: "", Rt: ""}, fmt.Errorf("error while set rt, %s", err)
		}
		return &pr.SignInResponse{At: at, Rt: rt}, nil
	}
	return &pr.SignInResponse{At: "", Rt: ""}, fmt.Errorf("error while loogin, wrong credentials %s", err)
}

func CheckPasswordHash(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}
