package services

import (
	"context"
	"github.com/995933447/apperrdef"
	"github.com/995933447/log-go"
	"github.com/995933447/optionstream"
	"github.com/995933447/reflectutil"
	"github.com/vision-first/wegod/internal/pkg/auth"
	"github.com/vision-first/wegod/internal/pkg/datamodel/models"
	"github.com/vision-first/wegod/internal/pkg/db/enum"
	"github.com/vision-first/wegod/internal/pkg/db/mysql/orms/gormimpl"
	"github.com/vision-first/wegod/internal/pkg/errs"
	"github.com/vision-first/wegod/internal/pkg/facades"
	"github.com/vision-first/wegod/internal/pkg/queryoptions"
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
		err = apperrdef.NewErr(errs.ErrCodeUserNotFound)
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

func (u *User) CreateUser(ctx context.Context, req *CreateUserReq) (*models.User, error) {
	userDO := &models.User{}
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

func (u *User) GetUserById(ctx context.Context, id uint64) (*models.User, error) {
	var userDO models.User
	if err := facades.MustGormDB(ctx, u.logger).First(&userDO, id).Error; err != nil {
		u.logger.Error(ctx, err)
		return nil, u.TransErr(err)
	}
	return &userDO, nil
}

func (u *User) EnsureNotUsePhoneUser(ctx context.Context, phone string) (bool, error) {
	var userNumForBoundReqPhone int64
	err := facades.MustGormDB(ctx, u.logger).Where(&models.User{Phone: phone}).Count(&userNumForBoundReqPhone).Error
	if err != nil {
		u.logger.Error(ctx, err)
		return false, u.TransErr(err)
	}
	return userNumForBoundReqPhone == 0, nil
}

func (u *User) GetUserByPhone(ctx context.Context, phone string) (*models.User, error) {
	var userDO models.User
	err := facades.MustGormDB(ctx, u.logger).Where(&models.User{Phone: phone}).First(ctx, &userDO).Error
	if err != nil {
		u.logger.Error(ctx, err)
		return nil, u.TransErr(err)
	}
	return &userDO, nil
}

func (u *User) PageUsers(ctx context.Context, queryStream *optionstream.QueryStream) ([]*models.User, *optionstream.Pagination, error) {
	db := facades.MustGormDB(ctx, u.logger)

	queryStreamProcessor := optionstream.NewQueryStreamProcessor(queryStream)
	queryStreamProcessor.
		OnUint64List(queryoptions.InIds, func(val []uint64) error {
			db.Where(enum.FieldId + " IN ?", val)
			return nil
		})

	var userDOs []*models.User
	pagination, err := queryStreamProcessor.PaginateFrom(ctx, gormimpl.NewOptStreamQuery(db), &userDOs)
	if err != nil {
		u.logger.Error(ctx, err)
		return nil, nil, u.TransErr(err)
	}

	return userDOs, pagination, nil
}