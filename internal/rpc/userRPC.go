package rpc

import (
	"context"
	"fmt"

	"github.com/IvanVojnic/bandEFuser/models"
	pr "github.com/IvanVojnic/bandEFuser/proto"

	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
)

type UserComm interface {
	GetFriends(ctx context.Context, userID uuid.UUID) (*[]models.User, error)
	SendFriendsRequest(ctx context.Context, userSender uuid.UUID, userReceiver uuid.UUID) error
	AcceptFriendsRequest(ctx context.Context, userSenderID uuid.UUID, userID uuid.UUID) error
	DeclineFriendsRequest(ctx context.Context, userSenderID uuid.UUID, userID uuid.UUID) error
	FindUser(ctx context.Context, userEmail string) (*models.User, error)
	GetRequest(ctx context.Context, userID uuid.UUID) (*[]models.User, error)
	GetUsers(ctx context.Context, usersID *[]uuid.UUID) (*[]models.User, error)
}

type UserCommServer struct {
	pr.UnimplementedUserCommServer
	userCommServ UserComm
}

func NewUserCommServer(userCommServ UserComm) *UserCommServer {
	return &UserCommServer{userCommServ: userCommServ}
}

func (s *UserCommServer) GetFriends(ctx context.Context, req *pr.GetFriendsRequest) (*pr.GetFriendsResponse, error) {
	userID, errParse := uuid.Parse(req.GetUserID())
	var users []*pr.User
	if errParse != nil {
		logrus.WithFields(logrus.Fields{
			"Error get user ID": errParse,
			"user":              userID,
		}).Errorf("error while parsing, %s", errParse)
		return &pr.GetFriendsResponse{Friends: users}, fmt.Errorf("error while parse userID, %s", errParse)
	}
	usersDB, err := s.userCommServ.GetFriends(ctx, userID)
	if err != nil {
		logrus.WithFields(logrus.Fields{
			"Error getting users friends": err,
			"user":                        usersDB,
		}).Errorf("error while getting users, %s", err)
		return &pr.GetFriendsResponse{Friends: users}, fmt.Errorf("error while gettingg friends from db, %s", err)
	}
	for _, user := range *usersDB {
		users = append(users, &pr.User{ID: user.ID.String(), Name: user.Name, Email: user.Email})
	}
	return &pr.GetFriendsResponse{Friends: users}, nil
}

func (s *UserCommServer) SendFriendsRequest(ctx context.Context, req *pr.SendFriendRequestReq) (*pr.SendFriendRequestResp, error) {
	userSenderID, errParse := uuid.Parse(req.GetUserID())
	if errParse != nil {
		logrus.WithFields(logrus.Fields{
			"Error get user ID": errParse,
			"user":              userSenderID,
		}).Errorf("error while parsing, %s", errParse)
		return &pr.SendFriendRequestResp{}, fmt.Errorf("error while parsing user sender ID, %s", errParse)
	}
	userReceiverID, errParse := uuid.Parse(req.GetReceiverID())
	if errParse != nil {
		logrus.WithFields(logrus.Fields{
			"Error get user ID": errParse,
			"user":              userReceiverID,
		}).Errorf("error while parsing, %s", errParse)
		return &pr.SendFriendRequestResp{}, fmt.Errorf("error while parsing user receiver ID, %s", errParse)
	}
	err := s.userCommServ.SendFriendsRequest(ctx, userSenderID, userReceiverID)
	if err != nil {
		logrus.WithFields(logrus.Fields{
			"Error sending request to be a friend": err,
		}).Errorf("error while inserting row into friends table, %s", errParse)
		return &pr.SendFriendRequestResp{}, fmt.Errorf("error while sending request, %s", err)
	}
	return &pr.SendFriendRequestResp{}, nil
}

func (s *UserCommServer) AcceptFriendsRequest(ctx context.Context, req *pr.AcceptFriendsRequestReq) (*pr.AcceptFriendsRequestResp, error) {
	userReceiverID, errParse := uuid.Parse(req.GetUserID())
	if errParse != nil {
		logrus.WithFields(logrus.Fields{
			"Error get user ID": errParse,
			"user":              userReceiverID,
		}).Errorf("error while parsing, %s", errParse)
		return &pr.AcceptFriendsRequestResp{}, fmt.Errorf("error while parsing user receiver ID, %s", errParse)
	}
	userSenderID, errParse := uuid.Parse(req.GetUserSenderID())
	if errParse != nil {
		logrus.WithFields(logrus.Fields{
			"Error get user ID": errParse,
			"user":              userSenderID,
		}).Errorf("error while parsing, %s", errParse)
		return &pr.AcceptFriendsRequestResp{}, fmt.Errorf("error while parsing user sender ID, %s", errParse)
	}
	err := s.userCommServ.AcceptFriendsRequest(ctx, userSenderID, userReceiverID)
	if err != nil {
		logrus.WithFields(logrus.Fields{
			"Error accepting request to be a friend": err,
		}).Errorf("error updating row into friends table, %s", err)
		return &pr.AcceptFriendsRequestResp{}, fmt.Errorf("error while accepting request, %s", err)
	}
	return &pr.AcceptFriendsRequestResp{}, nil
}

func (s *UserCommServer) DeclineFriendsRequest(ctx context.Context, req *pr.DeclineFriendsRequestReq) (*pr.DeclineFriendsRequestResp, error) {
	userReceiverID, errParse := uuid.Parse(req.UserID)
	if errParse != nil {
		logrus.WithFields(logrus.Fields{
			"Error get user ID": errParse,
			"user":              userReceiverID,
		}).Errorf("error while parsing, %s", errParse)
		return &pr.DeclineFriendsRequestResp{}, fmt.Errorf("error while parsing user receiver ID, %s", errParse)
	}
	userSenderID, errParse := uuid.Parse(req.UserSenderID)
	if errParse != nil {
		logrus.WithFields(logrus.Fields{
			"Error get user ID": errParse,
			"user":              userSenderID,
		}).Errorf("error while parsing, %s", errParse)
		return &pr.DeclineFriendsRequestResp{}, fmt.Errorf("error while parsing user sender ID, %s", errParse)
	}
	err := s.userCommServ.DeclineFriendsRequest(ctx, userSenderID, userReceiverID)
	if err != nil {
		logrus.WithFields(logrus.Fields{
			"Error updating request to not be a friend": err,
		}).Errorf("error insert row into friends table, %s", err)
		return &pr.DeclineFriendsRequestResp{}, fmt.Errorf("error while decling request, %s", err)
	}
	return &pr.DeclineFriendsRequestResp{}, nil
}

func (s *UserCommServer) FindUser(ctx context.Context, req *pr.FindUserRequest) (*pr.FindUserResponse, error) {
	user, err := s.userCommServ.FindUser(ctx, req.GetUserEmail())
	if err != nil {
		logrus.WithFields(logrus.Fields{
			"Error get user data": err,
			"user":                user,
		}).Errorf("error while parsing, %s", err)
		return &pr.FindUserResponse{}, fmt.Errorf("error while find user, %s", err)
	}
	return &pr.FindUserResponse{Friend: &pr.User{ID: user.ID.String(), Name: user.Name, Email: user.Email}}, nil
}

func (s *UserCommServer) GetRequest(ctx context.Context, req *pr.GetRequestReq) (*pr.GetRequestResp, error) {
	userID, errParse := uuid.Parse(req.GetUserID())
	var users []*pr.User
	if errParse != nil {
		logrus.WithFields(logrus.Fields{
			"Error get user ID": errParse,
			"user":              userID,
		}).Errorf("error while parsing, %s", errParse)
		return &pr.GetRequestResp{Users: users}, fmt.Errorf("error while parse userID, %s", errParse)
	}
	usersDB, err := s.userCommServ.GetRequest(ctx, userID)
	if err != nil {
		logrus.WithFields(logrus.Fields{
			"Error get users requests": err,
			"user":                     usersDB,
		}).Errorf("error while getting requests, %s", err)
		return &pr.GetRequestResp{Users: users}, fmt.Errorf("error while getting requests to be a friend from db, %s", err)
	}
	for _, user := range *usersDB {
		users = append(users, &pr.User{ID: user.ID.String(), Name: user.Name, Email: user.Email})
	}
	return &pr.GetRequestResp{Users: users}, nil
}

func (s *UserCommServer) GetUsers(ctx context.Context, req *pr.GetUsersRequest) (*pr.GetUsersResponse, error) {
	usersStrID := req.GetUsersID()
	var usersID []uuid.UUID
	for _, userStrID := range usersStrID {
		userID, errParseID := uuid.Parse(userStrID)
		if errParseID != nil {
			return &pr.GetUsersResponse{}, fmt.Errorf("error while parsing userID, %s", errParseID)
		}
		usersID = append(usersID, userID)
	}
	usersDB, err := s.userCommServ.GetUsers(ctx, &usersID)
	if err != nil {
		return &pr.GetUsersResponse{}, fmt.Errorf("error while getting users, %s", err)
	}
	var usersGRPC []*pr.User
	for _, user := range *usersDB {
		usersGRPC = append(usersGRPC, &pr.User{ID: user.ID.String(), Name: user.Name, Email: user.Email})
	}
	return &pr.GetUsersResponse{Users: usersGRPC}, nil
}
