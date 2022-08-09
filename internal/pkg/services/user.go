package services

import (
	"github.com/995933447/log-go"
	"github.com/vision-first/wegod/internal/pkg/datamodels"
	"golang.org/x/net/context"
)

type User struct {
	logger *log.Logger
}

type CreateUserReq struct {
	Phone string
	Password string
	NickName string
	Avatar string
	Desc string
	Gender uint8
}

func (u *User) CreateUser(ctx context.Context, req *CreateUserReq) (*datamodels.User, error) {
	//facades.MustGormDB(ctx, u.logger).Create(&datamodels.User{
	//
	//})
	return nil, nil
}