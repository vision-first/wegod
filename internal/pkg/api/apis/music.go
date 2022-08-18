package apis

import (
	"github.com/995933447/log-go"
	"github.com/vision-first/wegod/internal/pkg/api"
	"github.com/vision-first/wegod/internal/pkg/api/dtos"
)

type Music struct {
	logger log.Logger
}

func (m *Music) PageMusics(ctx api.Context, req *dtos.PageMusicsReq) (*dtos.PageMusicsResp, error) {
    var resp dtos.PageMusicsResp

    // TODO.Write your logic

    return &resp, nil
}