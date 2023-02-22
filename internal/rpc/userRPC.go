package rpc

import (
	"context"
	pr "userMS/proto"
)

type UserServer struct {
	pr.UserServer
}

func NewMs1Server() *UserServer {
	return &UserServer{}
}

func (s *UserServer) SignUp(ctx context.Context, req *pr.SignUpRequest) (*pr.SignUpResponse, error) {
	email := req.GetEmail()
	name := req.GetName()
	password := req.GetPassword()
}
