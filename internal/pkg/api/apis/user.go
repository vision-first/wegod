package apis

import (
	"github.com/995933447/log-go"
	"github.com/995933447/reflectutil"
	"github.com/vision-first/wegod/internal/pkg/api"
	"github.com/vision-first/wegod/internal/pkg/api/dtos"
	"github.com/vision-first/wegod/internal/pkg/facades"
	"github.com/vision-first/wegod/internal/pkg/services"
)

type User struct {
	logger *log.Logger
}

func NewUser(logger *log.Logger) *User {
	return &User{
		logger: logger,
	}
}

func (u *User) GetUserInfo(ctx api.Context, req *dtos.GetUserInfoReq) (*dtos.GetUserInfoResp, error) {
    var resp dtos.GetUserInfoResp

   	authIdent, err := ctx.GetAuthIdentOrFailed()
	if err != nil {
		u.logger.Error(ctx, err)
		return nil, err
	}

	userDO, err := services.NewUser(u.logger).GetUserById(ctx, authIdent.GetUserId())
	if err != nil {
		u.logger.Error(ctx, err)
		return nil, err
	}

	err = reflectutil.CopySameFields(userDO, &resp)
	if err != nil {
		u.logger.Error(ctx, err)
		return nil, err
	}

    return &resp, nil
}

func (u *User) SetUserInfo(ctx api.Context, req *dtos.SetUserInfoReq) (*dtos.SetUserInfoResp, error) {
    var resp dtos.SetUserInfoResp

	authIdent, err := ctx.GetAuthIdentOrFailed()
	if err != nil {
		u.logger.Error(ctx, err)
		return nil, err
	}

	userDO, err := services.NewUser(u.logger).GetUserById(ctx, authIdent.GetUserId())
	if err != nil {
		u.logger.Error(ctx, err)
		return nil, err
	}

	err = reflectutil.CopySameFields(req, userDO)
	if err != nil {
		u.logger.Error(ctx, err)
		return nil, err
	}

	err = facades.MustGormDB(ctx, u.logger).Save(userDO).Error
	if err != nil {
		u.logger.Error(ctx, err)
		return nil, err
	}

    return &resp, nil
}