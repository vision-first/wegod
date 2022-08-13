package services

import (
	"github.com/995933447/log-go"
	"github.com/995933447/reflectutil"
	"github.com/vision-first/wegod/internal/pkg/datamodels"
	"github.com/vision-first/wegod/internal/pkg/facades"
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
	userDO := &datamodels.User{}
	err := reflectutil.CopySameFields(req, userDO)
	if err != nil {
		return nil, err
	}

	if err = facades.MustGormDB(ctx, u.logger).Create(userDO).Error; err != nil {
		return nil, err
	}

	return userDO, nil
}

func (u *User) GetUserById(ctx context.Context, id uint64) (*datamodels.User, error) {
	var userDO datamodels.User
	if err := facades.MustGormDB(ctx, u.logger).First(&userDO, id).Error; err != nil {
		return nil, err
	}
	return &userDO, nil
}