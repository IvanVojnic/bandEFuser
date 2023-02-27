package rpc

import (
	"context"
	"fmt"
	pr "github.com/IvanVojnic/bandEFuser/proto"
	"github.com/google/uuid"
)

func (s *UserServer) GetFriends(ctx context.Context, req *pr.GetFriendsRequest) (*pr.GetFriendsResponse, error) {
	userID, errParse := uuid.Parse(req.GetUserID())
	var users []*pr.User
	if errParse != nil {
		return &pr.GetFriendsResponse{Friends: users}, fmt.Errorf("error while parse userID, %s", errParse)
	}
	err := s.userCommRepo.GetFriends(ctx, userID, users)
	if err != nil {
		return &pr.GetFriendsResponse{Friends: users}, fmt.Errorf("error while gettingg friends from db, %s", err)
	}
	return &pr.GetFriendsResponse{Friends: users}, nil
}

func (s *UserServer) SendFriendsRequest(ctx context.Context, req *pr.SendFriendRequestReq) (*pr.SendFriendRequestResp, error) {
	userSenderID, errParse := uuid.Parse(req.GetUserID())
	if errParse != nil {
		return &pr.SendFriendRequestResp{}, fmt.Errorf("error while parsing user sender ID, %s", errParse)
	}
	userReceiverID, errParse := uuid.Parse(req.GetReceiverID())
	if errParse != nil {
		return &pr.SendFriendRequestResp{}, fmt.Errorf("error while parsing user receiver ID, %s", errParse)
	}
	err := s.userCommRepo.SendFriendsRequest(ctx, userSenderID, userReceiverID)
	if err != nil {
		return &pr.SendFriendRequestResp{}, fmt.Errorf("error while sending request, %s", err)
	}
	return &pr.SendFriendRequestResp{}, nil
}

func (s *UserServer) AcceptFriendsRequest(ctx context.Context, req *pr.AcceptFriendsRequestReq) (*pr.AcceptFriendsRequestResp, error) {
	userReceiverID, errParse := uuid.Parse(req.GetUserID())
	if errParse != nil {
		return &pr.AcceptFriendsRequestResp{}, fmt.Errorf("error while parsing user receiver ID, %s", errParse)
	}
	userSenderID, errParse := uuid.Parse(req.GetUserSenderID())
	if errParse != nil {
		return &pr.AcceptFriendsRequestResp{}, fmt.Errorf("error while parsing user sender ID, %s", errParse)
	}
	err := s.userCommRepo.AcceptFriendsRequest(ctx, userSenderID, userReceiverID)
	if err != nil {
		return &pr.AcceptFriendsRequestResp{}, fmt.Errorf("error while accepting request, %s", err)
	}
	return &pr.AcceptFriendsRequestResp{}, nil
}

func (s *UserServer) DeclineFriendsRequest(ctx context.Context, req *pr.DeclineFriendsRequestReq) (*pr.DeclineFriendsRequestResp, error) {
	userReceiverID, errParse := uuid.Parse(req.UserID)
	if errParse != nil {
		return &pr.DeclineFriendsRequestResp{}, fmt.Errorf("error while parsing user receiver ID, %s", errParse)
	}
	userSenderID, errParse := uuid.Parse(req.UserSenderID)
	if errParse != nil {
		return &pr.DeclineFriendsRequestResp{}, fmt.Errorf("error while parsing user sender ID, %s", errParse)
	}
	err := s.userCommRepo.AcceptFriendsRequest(ctx, userSenderID, userReceiverID)
	if err != nil {
		return &pr.DeclineFriendsRequestResp{}, fmt.Errorf("error while decling request, %s", err)
	}
	return &pr.DeclineFriendsRequestResp{}, nil
}

func (s *UserServer) FindUser(ctx context.Context, req *pr.FindUserRequest) (*pr.FindUserResponse, error) {
	user, err := s.userCommRepo.FindUser(ctx, req.GetUserEmail())
	if err != nil {
		return &pr.FindUserResponse{}, fmt.Errorf("error while find user, %s", err)
	}
	return &pr.FindUserResponse{Friend: user}, nil
}

func (s *UserServer) GetRequest(ctx context.Context, req *pr.GetRequestReq) (*pr.GetRequestResp, error) {
	userID, errParse := uuid.Parse(req.GetUserID())
	var users []*pr.User
	if errParse != nil {
		return &pr.GetRequestResp{Users: users}, fmt.Errorf("error while parse userID, %s", errParse)
	}
	err := s.userCommRepo.GetRequest(ctx, userID, users)
	if err != nil {
		return &pr.GetRequestResp{Users: users}, fmt.Errorf("error while getting requests to be a friend from db, %s", err)
	}
	return &pr.GetRequestResp{Users: users}, nil
}
