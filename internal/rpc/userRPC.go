package rpc

import (
	"context"
	"fmt"

	"github.com/IvanVojnic/bandEFuser/models"
	pr "github.com/IvanVojnic/bandEFuser/proto"

	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
)

// UserComm interface define user comm methods implemented from service
type UserComm interface {
	GetFriends(ctx context.Context, userID uuid.UUID) ([]*models.User, error)
	SendFriendsRequest(ctx context.Context, userSender uuid.UUID, userReceiver uuid.UUID) error
	AcceptFriendsRequest(ctx context.Context, userSenderID uuid.UUID, userID uuid.UUID) error
	DeclineFriendsRequest(ctx context.Context, userSenderID uuid.UUID, userID uuid.UUID) error
	FindUser(ctx context.Context, userEmail string) (*models.User, error)
	GetRequest(ctx context.Context, userID uuid.UUID) ([]*models.User, error)
	GetUsers(ctx context.Context, usersID []*uuid.UUID) ([]*models.User, error)
	GetUser(ctx context.Context, userID uuid.UUID) (models.User, error)
}

// UserCommServer define user comm obj
type UserCommServer struct {
	pr.UnimplementedUserCommServer
	userCommServ UserComm
}

// NewUserCommServer used to init user comm odj
func NewUserCommServer(userCommServ UserComm) *UserCommServer {
	return &UserCommServer{userCommServ: userCommServ}
}

// GetFriends used to get friends
func (s *UserCommServer) GetFriends(ctx context.Context, req *pr.GetFriendsRequest) (*pr.GetFriendsResponse, error) { // nolint:dupl, gocritic
	userID, errParse := uuid.Parse(req.GetUserID())
	users := make([]*pr.User, 0)
	if errParse != nil {
		logrus.WithFields(logrus.Fields{
			"user ID": userID,
		}).Errorf("error while parsing, %s", errParse)
		return &pr.GetFriendsResponse{Friends: users}, fmt.Errorf("error while parse userID, %s", errParse)
	}
	usersDB, err := s.userCommServ.GetFriends(ctx, userID)
	if err != nil {
		logrus.WithFields(logrus.Fields{
			"users from database": usersDB,
		}).Errorf("error while getting users, %s", err)
		return &pr.GetFriendsResponse{Friends: users}, fmt.Errorf("error while gettingg friends from db, %s", err)
	}
	for _, user := range usersDB {
		users = append(users, &pr.User{ID: user.ID.String(), Name: user.Name, Email: user.Email})
	}
	return &pr.GetFriendsResponse{Friends: users}, nil
}

// SendFriendsRequest used to send request to a friend
func (s *UserCommServer) SendFriendsRequest(ctx context.Context, req *pr.SendFriendRequestReq) (*pr.SendFriendRequestResp, error) { // nolint:dupl, gocritic
	userSenderID, errParse := uuid.Parse(req.GetUserID())
	if errParse != nil {
		logrus.WithFields(logrus.Fields{
			"user sender ID": userSenderID,
		}).Errorf("error while parsing, %s", errParse)
		return &pr.SendFriendRequestResp{}, fmt.Errorf("error while parsing user sender ID, %s", errParse)
	}
	userReceiverID, errParse := uuid.Parse(req.GetReceiverID())
	if errParse != nil {
		logrus.WithFields(logrus.Fields{
			"user receiver ID": userReceiverID,
		}).Errorf("error while parsing, %s", errParse)
		return &pr.SendFriendRequestResp{}, fmt.Errorf("error while parsing user receiver ID, %s", errParse)
	}
	err := s.userCommServ.SendFriendsRequest(ctx, userSenderID, userReceiverID)
	if err != nil {
		logrus.Errorf("error while inserting row into friends table, %s", errParse)
		return &pr.SendFriendRequestResp{}, fmt.Errorf("error while sending request, %s", err)
	}
	return &pr.SendFriendRequestResp{}, nil
}

// AcceptFriendsRequest used to accept request to a friend
func (s *UserCommServer) AcceptFriendsRequest(ctx context.Context, req *pr.AcceptFriendsRequestReq) (*pr.AcceptFriendsRequestResp, error) { // nolint:dupl, gocritic
	userReceiverID, errParse := uuid.Parse(req.GetUserID())
	if errParse != nil {
		logrus.WithFields(logrus.Fields{
			"user receiver ID": userReceiverID,
		}).Errorf("error while parsing, %s", errParse)
		return &pr.AcceptFriendsRequestResp{}, fmt.Errorf("error while parsing user receiver ID, %s", errParse)
	}
	userSenderID, errParse := uuid.Parse(req.GetUserSenderID())
	if errParse != nil {
		logrus.WithFields(logrus.Fields{
			"user sender ID": userSenderID,
		}).Errorf("error while parsing, %s", errParse)
		return &pr.AcceptFriendsRequestResp{}, fmt.Errorf("error while parsing user sender ID, %s", errParse)
	}
	err := s.userCommServ.AcceptFriendsRequest(ctx, userSenderID, userReceiverID)
	if err != nil {
		logrus.Errorf("error updating row into friends table, %s", err)
		return &pr.AcceptFriendsRequestResp{}, fmt.Errorf("error while accepting request, %s", err)
	}
	return &pr.AcceptFriendsRequestResp{}, nil
}

// DeclineFriendsRequest used to decline request to a friend
func (s *UserCommServer) DeclineFriendsRequest(ctx context.Context, req *pr.DeclineFriendsRequestReq) (*pr.DeclineFriendsRequestResp, error) { // nolint:dupl, gocritic
	userReceiverID, errParse := uuid.Parse(req.UserID)
	if errParse != nil {
		logrus.WithFields(logrus.Fields{
			"user receiver ID": userReceiverID,
		}).Errorf("error while parsing, %s", errParse)
		return &pr.DeclineFriendsRequestResp{}, fmt.Errorf("error while parsing user receiver ID, %s", errParse)
	}
	userSenderID, errParse := uuid.Parse(req.UserSenderID)
	if errParse != nil {
		logrus.WithFields(logrus.Fields{
			"user sender ID": userSenderID,
		}).Errorf("error while parsing, %s", errParse)
		return &pr.DeclineFriendsRequestResp{}, fmt.Errorf("error while parsing user sender ID, %s", errParse)
	}
	err := s.userCommServ.DeclineFriendsRequest(ctx, userSenderID, userReceiverID)
	if err != nil {
		logrus.Errorf("error insert row into friends table, %s", err)
		return &pr.DeclineFriendsRequestResp{}, fmt.Errorf("error while decling request, %s", err)
	}
	return &pr.DeclineFriendsRequestResp{}, nil
}

// FindUser used to find user
func (s *UserCommServer) FindUser(ctx context.Context, req *pr.FindUserRequest) (*pr.FindUserResponse, error) {
	user, err := s.userCommServ.FindUser(ctx, req.GetUserEmail())
	if err != nil {
		logrus.WithFields(logrus.Fields{
			"user": user,
		}).Errorf("error while parsing, %s", err)
		return &pr.FindUserResponse{}, fmt.Errorf("error while find user, %s", err)
	}
	return &pr.FindUserResponse{Friend: &pr.User{ID: user.ID.String(), Name: user.Name, Email: user.Email}}, nil
}

// GetRequest used to get requests to be a friend
func (s *UserCommServer) GetRequest(ctx context.Context, req *pr.GetRequestReq) (*pr.GetRequestResp, error) { // nolint:dupl, gocritic
	userID, errParse := uuid.Parse(req.GetUserID())
	users := make([]*pr.User, 0)
	if errParse != nil {
		logrus.WithFields(logrus.Fields{
			"user ID": userID,
		}).Errorf("error while parsing (get request), %s", errParse)
		return &pr.GetRequestResp{Users: users}, fmt.Errorf("error while parse userID, %s", errParse)
	}
	usersDB, err := s.userCommServ.GetRequest(ctx, userID)
	if err != nil {
		logrus.WithFields(logrus.Fields{
			"users from database": usersDB,
		}).Errorf("error while getting requests, %s", err)
		return &pr.GetRequestResp{Users: users}, fmt.Errorf("error while getting requests to be a friend from db, %s", err)
	}
	for _, user := range usersDB {
		users = append(users, &pr.User{ID: user.ID.String(), Name: user.Name, Email: user.Email})
	}
	return &pr.GetRequestResp{Users: users}, nil
}

// GetUsers used to get users by theirs ID
func (s *UserCommServer) GetUsers(ctx context.Context, req *pr.GetUsersRequest) (*pr.GetUsersResponse, error) {
	usersStrID := req.GetUsersID()
	usersID := make([]*uuid.UUID, 0)
	for _, userStrID := range usersStrID {
		userID, errParseID := uuid.Parse(userStrID)
		if errParseID != nil {
			logrus.WithFields(logrus.Fields{
				"user ID": userID,
			}).Errorf("error while parsing (get request), %s", errParseID)
			return &pr.GetUsersResponse{}, fmt.Errorf("error while parsing userID, %s", errParseID)
		}
		usersID = append(usersID, &userID)
	}
	usersDB, err := s.userCommServ.GetUsers(ctx, usersID)
	if err != nil {
		logrus.WithFields(logrus.Fields{
			"users from database": usersDB,
		}).Errorf("error while getting requests, %s", err)
		return &pr.GetUsersResponse{}, fmt.Errorf("error while getting users, %s", err)
	}
	usersGRPC := make([]*pr.User, 0)
	for _, user := range usersDB {
		usersGRPC = append(usersGRPC, &pr.User{ID: user.ID.String(), Name: user.Name, Email: user.Email})
	}
	return &pr.GetUsersResponse{Users: usersGRPC}, nil
}

// GetUser used to get users by theirs ID
func (s *UserCommServer) GetUser(ctx context.Context, req *pr.GetUserRequest) (*pr.GetUserResponse, error) {
	userID, errParseID := uuid.Parse(req.GetUserID())
	if errParseID != nil {
		logrus.WithFields(logrus.Fields{
			"user ID": userID,
		}).Errorf("error while parsing, %s", errParseID)
		return &pr.GetUserResponse{}, fmt.Errorf("error while parsing userID, %s", errParseID)
	}
	userDB, err := s.userCommServ.GetUser(ctx, userID)
	if err != nil {
		logrus.WithFields(logrus.Fields{
			"user from database": userDB,
		}).Errorf("error while getting requests, %s", err)
		return &pr.GetUserResponse{}, fmt.Errorf("error while getting user, %s", err)
	}

	return &pr.GetUserResponse{User: &pr.User{ID: userDB.ID.String(), Name: userDB.Name, Email: userDB.Email})}, nil
}
