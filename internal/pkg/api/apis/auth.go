package apis

import (
	"github.com/995933447/log-go"
	"github.com/995933447/stringhelper-go"
	"github.com/vision-first/wegod/internal/pkg/api"
	"github.com/vision-first/wegod/internal/pkg/api/dtos"
	"github.com/vision-first/wegod/internal/pkg/enums"
	"github.com/vision-first/wegod/internal/pkg/services"
)

type Auth struct {
	logger log.Logger
}

func (a *Auth) Register(ctx api.Context, req *dtos.RegisterReq) (*dtos.RegisterResp, error) {
    var resp dtos.RegisterResp

    // TODO.Write your logic

    return &resp, nil
}

func (a *Auth) Login(ctx api.Context, req *dtos.LoginReq) (*dtos.LoginResp, error) {
    var resp dtos.LoginResp

    // TODO.Write your logic

    return &resp, nil
}

func (a *Auth) SendVerifyCodeForRegister(ctx api.Context, req *dtos.SendVerifyCodeForRegisterReq) (*dtos.SendVerifyCodeForRegisterResp, error) {
    var resp dtos.SendVerifyCodeForRegisterResp

    code := stringhelper.GenRandomStr(stringhelper.RandomStringModNumber, 6)
	err := services.NewSMS(a.logger).SendMsg(enums.MakeVerifyCodeMsgForRegister(code))
	if err != nil {
		a.logger.Error(ctx, err)
		return nil, err
	}

    return &resp, nil
}