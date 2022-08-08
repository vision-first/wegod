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
		buddhaDataModels []*datamodels.Buddha
	)
	buddhaDataModels, resp.Pagination, err = services.NewBuddha(b.logger).
		PageBuddha(ctx, optionstream.NewQueryStream(nil, req.Limit, req.Offset))
	if err != nil {
		b.logger.Error(ctx, err)
		return nil, err
	}
	for _, buddhaDataModel := range buddhaDataModels {
		resp.List = append(resp.List, &dtos.Buddha{
			Id: buddhaDataModel.Id,
			Name: buddhaDataModel.Name,
			Image: buddhaDataModel.Image,
			Sort: buddhaDataModel.Sort,
		})
	}
	return &resp, nil
}

func (b *Buddha) WatchBuddha(ctx api.Context, req *dtos.WatchBuddhaReq) (*dtos.WatchBuddhaResp, error) {
    var resp dtos.WatchBuddhaResp

    // TODO.Write your logic

    return &resp, nil
}

func (b *Buddha) UnwatchBuddha(ctx api.Context, req *dtos.UnwatchBuddhaReq) (*dtos.UnwatchBuddhaResp, error) {
    var resp dtos.UnwatchBuddhaResp

    // TODO.Write your logic

    return &resp, nil
}