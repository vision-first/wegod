package services

import (
	"github.com/995933447/apperrdef"
	"github.com/995933447/log-go"
	"github.com/995933447/reflectutil"
	"github.com/vision-first/wegod/internal/pkg/auth"
	"github.com/vision-first/wegod/internal/pkg/datamodels"
	"github.com/vision-first/wegod/internal/pkg/errs"
	"github.com/vision-first/wegod/internal/pkg/facades"
	"golang.org/x/net/context"
	"gorm.io/gorm"
)

type User struct {
	logger *log.Logger
}

func NewUser(logger *log.Logger) *User {
	return &User{
		logger: logger,
	}
}

func (*User) TransErr(err error) error {
	switch err {
	case gorm.ErrRecordNotFound:
		err = apperrdef.NewErr(errs.ErrUserNotFound)
	}
	return err
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
		u.logger.Error(ctx, err)
		return nil, u.TransErr(err)
	}

	userDO.HashedPassword, err = auth.HashPassword(req.Password)
	if err != nil {
		u.logger.Error(ctx, err)
		return nil, u.TransErr(err)
	}

	if err = facades.MustGormDB(ctx, u.logger).Create(userDO).Error; err != nil {
		u.logger.Error(ctx, err)
		return nil, u.TransErr(err)
	}

	return userDO, nil
}

func (u *User) GetUserById(ctx context.Context, id uint64) (*datamodels.User, error) {
	var userDO datamodels.User
	if err := facades.MustGormDB(ctx, u.logger).First(&userDO, id).Error; err != nil {
		u.logger.Error(ctx, err)
		return nil, u.TransErr(err)
	}
	return &userDO, nil
}

func (u *User) EnsureNotUsePhoneUser(ctx context.Context, phone string) (bool, error) {
	var userNumForBoundReqPhone int64
	err := facades.MustGormDB(ctx, u.logger).Where(&datamodels.User{Phone: phone}).Count(&userNumForBoundReqPhone).Error
	if err != nil {
		u.logger.Error(ctx, err)
		return false, u.TransErr(err)
	}
	return userNumForBoundReqPhone == 0, nil
}

func (u *User) GetUserByPhone(ctx context.Context, phone string) (*datamodels.User, error) {
	var userDO datamodels.User
	err := facades.MustGormDB(ctx, u.logger).Where(&datamodels.User{Phone: phone}).First(ctx, &userDO).Error
	if err != nil {
		u.logger.Error(ctx, err)
		return nil, u.TransErr(err)
	}
	return &userDO, nil
}