package apis

import (
	"github.com/995933447/log-go"
	"github.com/vision-first/wegod/internal/pkg/api"
	"github.com/vision-first/wegod/internal/pkg/api/dtos"
)

type Donation struct {
	logger log.Logger
}

func (d *Donation) CreateDonationOrder(ctx api.Context, req *dtos.CreateDonationOrderReq) (*dtos.CreateDonationOrderResp, error) {
    var resp dtos.CreateDonationOrderResp

    // TODO.Write your logic

    return &resp, nil
}