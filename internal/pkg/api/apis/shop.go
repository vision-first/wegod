package apis

import (
	"github.com/995933447/log-go"
	"github.com/vision-first/wegod/internal/pkg/api"
	"github.com/vision-first/wegod/internal/pkg/api/dtos"
)

type Shop struct {
	logger log.Logger
}

func (s *Shop) PageProducts(ctx api.Context, req *dtos.PageProductsReq) (*dtos.PageProductsResp, error) {
    var resp dtos.PageProductsResp

    // TODO.Write your logic

    return &resp, nil
}