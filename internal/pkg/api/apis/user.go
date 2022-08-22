package apis

import (
	"github.com/995933447/log-go"
	"github.com/vision-first/wegod/internal/pkg/api"
	"github.com/vision-first/wegod/internal/pkg/api/dtos"
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

    // TODO.write your logic

    return &resp, nil
}

func (u *User) SetUserInfo(ctx api.Context, req *dtos.SetUserInfoReq) (*dtos.SetUserInfoResp, error) {
    var resp dtos.SetUserInfoResp

    // TODO.write your logic

    return &resp, nil
}