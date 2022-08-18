package apis

import (
	"github.com/995933447/log-go"
	"github.com/995933447/optionstream"
	"github.com/995933447/reflectutil"
	"github.com/vision-first/wegod/internal/pkg/api"
	"github.com/vision-first/wegod/internal/pkg/api/dtos"
	"github.com/vision-first/wegod/internal/pkg/datamodel/metadata"
	"github.com/vision-first/wegod/internal/pkg/datamodel/models"
	"github.com/vision-first/wegod/internal/pkg/queryoptions"
	"github.com/vision-first/wegod/internal/pkg/services"
)

type PrayProp struct {
	logger *log.Logger
}

func (p *PrayProp) CreatePrayPropOrder(ctx api.Context, req *dtos.CreatePrayPropOrderReq) (*dtos.CreatePrayPropOrderResp, error) {
    var resp dtos.CreatePrayPropOrderResp

	authIdent, err := ctx.GetAuthIdentOrFailed()
	if err != nil {
		p.logger.Error(ctx, err)
		return nil, err
	}

	_, err = services.NewPrayPropOrder(p.logger).CreateOrder(ctx, authIdent.GetUserId(), req.PrayPropId, req.Num)
	if err != nil {
		p.logger.Error(ctx, err)
		return nil, err
	}

    return &resp, nil
}

func (p *PrayProp) PagePrayProps(ctx api.Context, req *dtos.PagePrayPropsReq) (*dtos.PagePrayPropsResp, error) {
    var resp dtos.PagePrayPropsResp

	queryStream := optionstream.NewQueryStream(req.QueryOptions, req.Limit, req.Offset)
	queryStream.SetOption(queryoptions.OnShelfStatus, nil)

    var (
    	propDOs []*models.PrayProp
    	err error
	)
	propDOs, resp.Pagination, err = services.NewPrayProp(p.logger).PagePrayProps(ctx, queryStream)
	if err != nil {
		p.logger.Error(ctx, err)
		return nil, err
	}

	resp.List = transPrayPropDOsToDTOs(propDOs)

    return &resp, nil
}

func (p *PrayProp) PageUserPrayProps(ctx api.Context, req *dtos.PageUserPrayPropsReq) (*dtos.PageUserPrayPropsResp, error) {
	var (
		resp dtos.PageUserPrayPropsResp
		userPrayPropDOs []*models.UserPrayProp
		err error
	)
	userPrayPropDOs, resp.Pagination, err = services.NewUserPrayProp(p.logger).
		PageUserPrayProps(ctx, optionstream.NewQueryStream(req.QueryOptions, req.Limit, req.Offset))
	if err != nil {
		p.logger.Error(ctx, err)
		return nil, err
	}

	queryPropStream := optionstream.NewQueryStream(
		[]*optionstream.Option{
			{Key: queryoptions.InIds, Val: reflectutil.PluckUint64(userPrayPropDOs, metadata.FieldPropId)},
		}, 0, 0)
	propDOs, _, err := services.NewPrayProp(p.logger).PagePrayProps(ctx, queryPropStream)
	if err != nil {
		p.logger.Error(ctx, err)
		return nil, err
	}

	propDOMap := reflectutil.MapByKey(propDOs, metadata.FieldId).(map[uint64]*models.PrayProp)
	for _, userPrayPropDO := range userPrayPropDOs {
		propDO, ok := propDOMap[userPrayPropDO.PropId]
		if !ok {
			continue
		}
		userPrayDTO := &dtos.UserPrayProp{}
		userPrayDTO.PrayProp = *transPrayPropDOToDTO(propDO)
		userPrayDTO.Num = userPrayPropDO.Num
		resp.List = append(resp.List, userPrayDTO)
	}

	return &resp, nil
}
