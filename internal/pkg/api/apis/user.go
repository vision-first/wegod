package apis

import (
	"github.com/995933447/log-go"
	"github.com/vision-first/wegod/internal/pkg/api"
	"github.com/vision-first/wegod/internal/pkg/api/dtos"
)

type User struct {
	logger log.Logger
}

func (u *User) GetUser(ctx api.Context, req *dtos.GetUserReq) (*dtos.GetUserResp, error) {
    var resp dtos.GetUserResp

    // TODO.Write your logic

    return &resp, nil
}

