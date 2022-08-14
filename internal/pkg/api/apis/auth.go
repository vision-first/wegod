package apis

import (
	"github.com/995933447/apperrdef"
	"github.com/995933447/log-go"
	"github.com/995933447/redisgroup"
	"github.com/995933447/stringhelper-go"
	"github.com/vision-first/wegod/internal/pkg/api"
	"github.com/vision-first/wegod/internal/pkg/api/dtos"
	"github.com/vision-first/wegod/internal/pkg/auth"
	"github.com/vision-first/wegod/internal/pkg/enums"
	"github.com/vision-first/wegod/internal/pkg/errs"
	"github.com/vision-first/wegod/internal/pkg/services"
)

type Auth struct {
	logger *log.Logger
	redisGroup redisgroup.Group
}

func (a *Auth) TransErr(err error) error {
	return err
}

func (a *Auth) Register(ctx api.Context, req *dtos.RegisterReq) (*dtos.RegisterResp, error) {
    var resp dtos.RegisterResp

	userSrv := services.NewUser(a.logger)
	notExist, err := userSrv.EnsureNotUsePhoneUser(ctx, req.Phone)
	if err != nil {
		a.logger.Error(ctx, err)
		return nil, a.TransErr(err)
	}

	if !notExist {
		err = apperrdef.NewErr(errs.ErrPhoneBeenRegisteredByOtherUser)
		a.logger.Error(ctx, err)
		return nil, a.TransErr(err)
	}

	authSrv := services.NewAuth(a.logger)
	ok, err := authSrv.AuthPhoneVerifyCode(ctx, req.Phone, req.VerifyCode)
	if err != nil {
		a.logger.Error(ctx, err)
		return nil, a.TransErr(err)
	}

	if !ok {
		err = apperrdef.NewErr(errs.ErrBadPhoneVerifyCode)
		a.logger.Error(ctx, err)
		return nil, a.TransErr(err)
	}

	userDO, err := userSrv.CreateUser(ctx, &services.CreateUserReq{
		Phone: req.Phone,
		Password: req.Password,
	})
	if err != nil {
		a.logger.Error(ctx, err)
		return nil, a.TransErr(err)
	}

	resp.Token, err = authSrv.CreateAuthTokenForUser(ctx, userDO.Id)
	if err != nil {
		a.logger.Error(ctx, err)
		return nil, a.TransErr(err)
	}

    return &resp, nil
}

func (a *Auth) Login(ctx api.Context, req *dtos.LoginReq) (*dtos.LoginResp, error) {
    var resp dtos.LoginResp

    userSrv := services.NewUser(a.logger)
	userDO, err := userSrv.GetUserByPhone(ctx, req.Phone)
	if err != nil {
		a.logger.Error(ctx, err)
		return nil, a.TransErr(err)
	}

	authSrv := services.NewAuth(a.logger)

	if req.IsByCode {
		ok, err := authSrv.AuthPhoneVerifyCode(ctx, req.Phone, req.Code)
		if err != nil {
			a.logger.Error(ctx, err)
			return nil, a.TransErr(err)
		}

		if !ok {
			err = apperrdef.NewErr(errs.ErrBadPhoneVerifyCode)
			a.logger.Error(ctx, err)
			return nil, a.TransErr(err)
		}
	} else {
		ok, err := auth.EqualPassword(req.Password, userDO.HashedPassword)
		if err != nil {
			a.logger.Error(ctx, err)
			return nil, a.TransErr(err)
		}

		if !ok {
			err = apperrdef.NewErr(errs.ErrPasswordNotCorrect)
			a.logger.Error(ctx, err)
			return nil, a.TransErr(err)
		}
	}

	resp.Token, err = authSrv.CreateAuthTokenForUser(ctx, userDO.Id)
	if err != nil {
		a.logger.Error(ctx, err)
		return nil, a.TransErr(err)
	}

    return &resp, nil
}

func (a *Auth) SendVerifyCodeForRegister(ctx api.Context, req *dtos.SendVerifyCodeForRegisterReq) (*dtos.SendVerifyCodeForRegisterResp, error) {
    var resp dtos.SendVerifyCodeForRegisterResp

	notExist, err := services.NewUser(a.logger).EnsureNotUsePhoneUser(ctx, req.Phone)
	if err != nil {
		a.logger.Error(ctx, err)
		return nil, a.TransErr(err)
	}

	if !notExist {
		err = apperrdef.NewErr(errs.ErrPhoneBeenRegisteredByOtherUser)
		a.logger.Error(ctx, err)
		return nil, a.TransErr(err)
	}

    code := stringhelper.GenRandomStr(stringhelper.RandomStringModNumber, 6)
	err = services.NewSMS(a.logger).SendMsg(req.Phone, enums.MakeVerifyCodeMsgForRegister(code))
	if err != nil {
		a.logger.Error(ctx, err)
		return nil, a.TransErr(err)
	}

	err = services.NewAuth(a.logger).RememberPhoneVerifyCodeForAuth(ctx, req.Phone, code)
	if err != nil {
		a.logger.Error(ctx, err)
		return nil, a.TransErr(err)
	}

    return &resp, nil
}

func (a *Auth) SendVerifyCodeForLogin(ctx api.Context, req *dtos.SendVerifyCodeForLoginReq) (*dtos.SendVerifyCodeForLoginResp, error) {
    var resp dtos.SendVerifyCodeForLoginResp

	code := stringhelper.GenRandomStr(stringhelper.RandomStringModNumber, 6)
	err := services.NewSMS(a.logger).SendMsg(req.Phone, enums.MakeVerifyCodeMsgForRegister(code))
	if err != nil {
		a.logger.Error(ctx, err)
		return nil, a.TransErr(err)
	}

	err = services.NewAuth(a.logger).RememberPhoneVerifyCodeForAuth(ctx, req.Phone, code)
	if err != nil {
		a.logger.Error(ctx, err)
		return nil, a.TransErr(err)
	}

    return &resp, nil
}