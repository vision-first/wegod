package apis

import (
	"github.com/995933447/eventobserver"
	"github.com/995933447/log-go"
	"github.com/995933447/optionstream"
	"github.com/995933447/reflectutil"
	"github.com/vision-first/wegod/internal/pkg/api"
	"github.com/vision-first/wegod/internal/pkg/api/dtos"
	"github.com/vision-first/wegod/internal/pkg/datamodel/metadata"
	"github.com/vision-first/wegod/internal/pkg/datamodel/models"
	"github.com/vision-first/wegod/internal/pkg/event"
	"github.com/vision-first/wegod/internal/pkg/event/payloads"
	"github.com/vision-first/wegod/internal/pkg/facades"
	"github.com/vision-first/wegod/internal/pkg/queryoptions"
	"github.com/vision-first/wegod/internal/pkg/services"
)

type Donation struct {
	logger *log.Logger
}

func (d *Donation) CreateDonationOrder(ctx api.Context, req *dtos.CreateDonationOrderReq) (*dtos.CreateDonationOrderResp, error) {
    var resp dtos.CreateDonationOrderResp

	authIdent, err := ctx.GetAuthIdentOrFailed()
	if err != nil {
		d.logger.Error(ctx, err)
		return nil, err
	}

    order, err := services.NewDonationOrder(d.logger).CreateOrder(ctx, &services.CreateOrderReq{
		UserId: authIdent.GetUserId(),
		BuddhaId: req.BuddhaId,
		Money: req.Money,
		Scene: models.DonationScene(req.Scene),
	})
	if err != nil {
		d.logger.Error(ctx, err)
		return nil, err
	}

	facades.EventDispatcher(d.logger).
		Dispatch(ctx, *eventobserver.NewEvent(event.EventNameCreatedDonationOrder, &payloads.CreatedDonationOrder{
			UserId: authIdent.GetUserId(),
			Money: req.Money,
			Id: order.Id,
			CreatedAt: order.CreatedAt,
		}))

    return &resp, nil
}

func (d *Donation) PageDonationRanks(ctx api.Context, req *dtos.PageDonationRanksReq) (*dtos.PageDonationRanksResp, error) {
	queryStatStream := optionstream.NewQueryStream(req.QueryOptions, req.Limit, req.Offset)
	queryStatStream.SetOption(queryoptions.OrderByPayedMoneyDesc, nil)

	var (
		resp dtos.PageDonationRanksResp
		statDOs []*models.UserDonationDailyStat
		err error
	)
	statDOs, resp.Pagination, err = services.NewUserDonationStat(d.logger).PageUserDonationRanks(ctx, queryStatStream)
	if err != nil {
		d.logger.Error(ctx, err)
		return nil, err
	}

	if len(statDOs) == 0 {
		return &resp, nil
	}

	userDOs, _, err := services.NewUser(d.logger).PageUsers(
		ctx,
		optionstream.NewQueryStream(
			[]*optionstream.Option{
				{Key: queryoptions.InIds, Val: reflectutil.PluckUint64(statDOs, metadata.FieldUserId)},
			},
			0,
			0,
			),
		)
	if err != nil {
		d.logger.Error(ctx, err)
		return nil, err
	}

	userMap := reflectutil.MapByKey(userDOs, metadata.FieldId).(map[uint64]*models.User)
	for _, statDO := range statDOs {
		userDO, ok := userMap[statDO.UserId]
		if !ok {
			continue
		}
		resp.List = append(resp.List, &dtos.UserDonationRecord{
			Money: statDO.TotalPayedMoney,
			NickName: userDO.NickName,
		})
	}

	return &resp, nil
}