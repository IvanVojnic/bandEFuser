package main

import (
	"github.com/IvanVojnic/bandEFuser/internal/config"
	"github.com/IvanVojnic/bandEFuser/internal/repository"
	"github.com/sirupsen/logrus"
)

func main() {
	cfg, err := config.NewConfig()
	if err != nil {
		logrus.WithFields(logrus.Fields{
			"Error":  err,
			"config": cfg,
		}).Fatal("failed to get config")
	}
	var profileServ *service.AuthService
	var userServ *service.UserCommSrv
	db, err := repository.NewPostgresDB(cfg)
	if err != nil {
		logrus.WithFields(logrus.Fields{
			"Error connection to database rep.NewPostgresDB()": err,
		}).Fatal("DB ERROR CONNECTION")
	}
	defer repository.ClosePool(db)
	profileRepo := repository.NewUserPostgres(db)
	userRepo := repository.NewUserCommPostgres(db)
	profileServ = service.NewAuthService(profileRepo)
	userServ = service.NewUserCommSrv(userRepo)
}
