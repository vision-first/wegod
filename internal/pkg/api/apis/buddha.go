package apis

import (
	"github.com/995933447/log-go"
	"github.com/995933447/optionstream"
	"github.com/vision-first/wegod/internal/pkg/api"
	"github.com/vision-first/wegod/internal/pkg/api/dtos"
	"github.com/vision-first/wegod/internal/pkg/datamodels"
	"github.com/vision-first/wegod/internal/pkg/services"
)

type Buddha struct {
	logger *log.Logger
}

func NewBuddha(logger *log.Logger) *Buddha {
	return &Buddha{
		logger: logger,
	}
}

func (b *Buddha) PageBuddha(ctx api.Context, req *dtos.PageBuddhaReq) (*dtos.PageBuddhaResp, error) {
	var (
		resp dtos.PageBuddhaResp
		err error
		buddhaDOs []*datamodels.Buddha
	)
	buddhaDOs, resp.Pagination, err = services.NewBuddha(b.logger).
		PageBuddha(ctx, optionstream.NewQueryStream(nil, req.Limit, req.Offset))
	if err != nil {
		b.logger.Error(ctx, err)
		return nil, err
	}

	for _, buddhaDO := range buddhaDOs {
		resp.List = append(resp.List, &dtos.Buddha{
			Id: buddhaDO.Id,
			Name: buddhaDO.Name,
			Image: buddhaDO.Image,
			Sort: buddhaDO.Sort,
		})
	}

	return &resp, nil
}

func (b *Buddha) WatchBuddha(ctx api.Context, req *dtos.WatchBuddhaReq) (*dtos.WatchBuddhaResp, error) {
    var resp dtos.WatchBuddhaResp

	authIdent, err := ctx.GetAuthIdentOrFailed()
	if err != nil {
		b.logger.Error(ctx, err)
		return nil, err
	}

	if err = services.NewBuddha(b.logger).WatchBuddha(ctx, authIdent.GetUserId(), req.BuddhaId); err != nil {
		b.logger.Error(ctx, err)
		return nil, err
	}

    return &resp, nil
}

func (b *Buddha) UnwatchBuddha(ctx api.Context, req *dtos.UnwatchBuddhaReq) (*dtos.UnwatchBuddhaResp, error) {
    var resp dtos.UnwatchBuddhaResp

	authIdent, err := ctx.GetAuthIdentOrFailed()
	if err != nil {
		b.logger.Error(ctx, err)
		return nil, err
	}

    err = services.NewBuddha(b.logger).UnwatchBuddha(ctx, authIdent.GetUserId(), req.BuddhaId)
	if err != nil {
		b.logger.Error(ctx, err)
		return nil, err
	}

    return &resp, nil
}


func (b *Buddha) PageUserWatchedBuddha(ctx api.Context, req *dtos.PageUserWatchedBuddhaReq) (*dtos.PageUserWatchedBuddhaResp, error) {
	var (
		resp dtos.PageUserWatchedBuddhaResp
		buddhaDOs []*datamodels.Buddha
		err error
	)
	buddhaDOs, resp.Pagination, err = services.NewBuddha(b.logger).
		PageUserWatchedBuddha(ctx, optionstream.NewQueryStream(req.QueryOptions, req.Limit, req.Offset))
	if err != nil {
		b.logger.Error(ctx, err)
		return nil, err
	}

	for _, buddhaDO := range buddhaDOs {
		resp.List = append(resp.List, &dtos.Buddha{
			Id: buddhaDO.Id,
			Name: buddhaDO.Name,
			Image: buddhaDO.Image,
			Sort: buddhaDO.Sort,
		})
	}

    return &resp, nil
}