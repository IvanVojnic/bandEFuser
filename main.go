// Package main is a main package used to start the program
package main

import (
	"fmt"
	"net"

	"github.com/IvanVojnic/bandEFuser/internal/config"
	"github.com/IvanVojnic/bandEFuser/internal/repository"
	"github.com/IvanVojnic/bandEFuser/internal/rpc"
	"github.com/IvanVojnic/bandEFuser/internal/service"
	pr "github.com/IvanVojnic/bandEFuser/proto"

	"github.com/sirupsen/logrus"

	"google.golang.org/grpc"
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
	db, err := repository.NewPostgresDB(cfg)
	if err != nil {
		logrus.WithFields(logrus.Fields{
			"Error connection to database rep.NewPostgresDB()": err,
		}).Fatal("DB ERROR CONNECTION")
	}
	defer repository.ClosePool(db)
	userAuthRepo := repository.NewUserPostgres(db)
	userCommRepo := repository.NewUserCommPostgres(db)

	userAuthServ := service.NewUserAuthServer(userAuthRepo)
	userCommServ := service.NewUserCommServer(userCommRepo)

	userAuthGRPC := rpc.NewUserAuthServer(userAuthServ)
	userCommGRPC := rpc.NewUserCommServer(userCommServ)

	pr.RegisterUserAuthServer(s, userAuthGRPC)
	pr.RegisterUserCommServer(s, userCommGRPC)
	listen, err := net.Listen("tcp", fmt.Sprintf("%s:%s", cfg.Host, cfg.Port))
	if err != nil {
		logrus.Fatalf("error while listening port: %v", err)
	}

	logrus.Info("a")

	if errServ := s.Serve(listen); errServ != nil {
		logrus.Fatalf("error while listening server: %v", err)
	}
}
