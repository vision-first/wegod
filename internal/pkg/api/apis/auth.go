package apis

import (
	"github.com/995933447/log-go"
	"github.com/vision-first/wegod/internal/pkg/api"
	"github.com/vision-first/wegod/internal/pkg/api/dtos"
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