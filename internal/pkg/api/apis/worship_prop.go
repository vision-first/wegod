package apis

import (
	"github.com/995933447/log-go"
	"github.com/995933447/optionstream"
	"github.com/vision-first/wegod/internal/pkg/api"
	"github.com/vision-first/wegod/internal/pkg/api/dtos"
	"github.com/vision-first/wegod/internal/pkg/datamodel/models"
	"github.com/vision-first/wegod/internal/pkg/queryoptions"
	"github.com/vision-first/wegod/internal/pkg/services"
)

type WorshipProp struct {
	logger *log.Logger
}

func (w *WorshipProp) PageWorshipProps(ctx api.Context, req *dtos.PageQueryReq) (*dtos.PageWorshipPropsResp, error) {
    var resp dtos.PageWorshipPropsResp

	authIdent, err := ctx.GetAuthIdentOrFailed()
	if err != nil {
		w.logger.Error(ctx, err)
		return nil, err
	}

	queryStream := optionstream.NewQueryStream(req.QueryOptions, req.Limit, req.Offset)
	queryStream.SetOption(queryoptions.EqualUserId, authIdent.GetUserId())

	var propDOs []*models.WorshipProp
	propDOs, resp.Pagination, err = services.NewWorshipProp(w.logger).PageWorshipProps(ctx, queryStream)
	if err != nil {
		w.logger.Error(ctx, err)
		return nil, err
	}

	resp.List = transWishPropDOsToDTOs(propDOs)

    return &resp, nil
}

func (w *WorshipProp) CreateWorshipPropOrder(ctx api.Context, req *dtos.CreateWorshipPropOrderReq) (*dtos.CreateWorshipPropOrderResp, error) {
    var resp dtos.CreateWorshipPropOrderResp

	authIdent, err := ctx.GetAuthIdentOrFailed()
	if err != nil {
		w.logger.Error(ctx, err)
		return nil, err
	}

    _, err = services.NewPrayPropOrder(w.logger).CreateOrder(ctx, authIdent.GetUserId(), req.WorshipPropId, req.Num)
	if err != nil {
		w.logger.Error(ctx, err)
		return nil, err
	}

    return &resp, nil
}