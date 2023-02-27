package main

import (
	"github.com/IvanVojnic/bandEFuser/internal/config"
	"github.com/IvanVojnic/bandEFuser/internal/repository"
	"github.com/IvanVojnic/bandEFuser/internal/rpc"
	pr "github.com/IvanVojnic/bandEFuser/proto"
	"github.com/sirupsen/logrus"
	"google.golang.org/grpc"
	"net"
)

func main() {

	var opts []grpc.ServerOption
	s := grpc.NewServer(opts...)

	cfg, err := config.NewConfig()
	if err != nil {
		logrus.WithFields(logrus.Fields{
			"Error":  err,
			"config": cfg,
		}).Fatal("failed to get config")
	}
	//var userAuthServ *rpc.UserServer
	//var userServ *rpc.UserCommServer
	db, err := repository.NewPostgresDB(cfg)
	if err != nil {
		logrus.WithFields(logrus.Fields{
			"Error connection to database rep.NewPostgresDB()": err,
		}).Fatal("DB ERROR CONNECTION")
	}
	defer repository.ClosePool(db)
	userAuthRepo := repository.NewUserPostgres(db)
	userCommRepo := repository.NewUserCommPostgres(db)
	userServ := rpc.NewUserAuthServer(userAuthRepo, userCommRepo)

	pr.RegisterUserServer(s, userServ)
	listen, err := net.Listen("tcp", ":8000")
	if err != nil {
		defer logrus.Fatalf("error while listening port: %e", err)
	}

	if errServ := s.Serve(listen); errServ != nil {
		defer logrus.Fatalf("error while listening server: %e", err)
	}
}
