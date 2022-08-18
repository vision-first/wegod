package apis

import (
	"github.com/995933447/apperrdef"
	"github.com/995933447/log-go"
	"github.com/995933447/optionstream"
	"github.com/vision-first/wegod/internal/pkg/api"
	"github.com/vision-first/wegod/internal/pkg/api/dtos"
	"github.com/vision-first/wegod/internal/pkg/datamodel/models"
	"github.com/vision-first/wegod/internal/pkg/db/enum"
	"github.com/vision-first/wegod/internal/pkg/errs"
	"github.com/vision-first/wegod/internal/pkg/queryoptions"
	"github.com/vision-first/wegod/internal/pkg/services"
	"time"
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
		buddhaDOs []*models.Buddha
	)
	buddhaDOs, resp.Pagination, err = services.NewBuddha(b.logger).
		PageBuddha(ctx, optionstream.NewQueryStream(nil, req.Limit, req.Offset))
	if err != nil {
		b.logger.Error(ctx, err)
		return nil, err
	}

	resp.List = transBuddhaDOsToDTOs(buddhaDOs)

	return &resp, nil
}

func (b *Buddha) WatchBuddha(ctx api.Context, req *dtos.WatchBuddhaReq) (*dtos.WatchBuddhaResp, error) {
    var resp dtos.WatchBuddhaResp

	buddhaSrv := services.NewBuddha(b.logger)
	isFreeRent, err := buddhaSrv.IsBuddhaFreeRent(ctx, req.BuddhaId)
	if err != nil {
		b.logger.Error(ctx, err)
		return nil, err
	}

	authIdent, err := ctx.GetAuthIdentOrFailed()
	if err != nil {
		b.logger.Error(ctx, err)
		return nil, err
	}

	var expireAt int64
	if !isFreeRent {
		expireAt, err = services.NewUserBuddhaRent(b.logger).GetRendBuddhaExpireAt(ctx, authIdent.GetUserId(), req.BuddhaId)
		if err != nil {
			b.logger.Error(ctx, err)
			return nil, err
		}

		if expireAt <= time.Now().Unix() {
			err = apperrdef.NewErr(errs.ErrCodeBuddhaRentExpired)
			b.logger.Error(ctx, err)
			return nil, err
		}
	}

	if err = buddhaSrv.WatchBuddha(ctx, authIdent.GetUserId(), req.BuddhaId, expireAt); err != nil {
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

func (b *Buddha) PageUserWatchedBuddhas(ctx api.Context, req *dtos.PageUserWatchedBuddhasReq) (*dtos.PageUserWatchedBuddhasResp, error) {
	authIdent, err := ctx.GetAuthIdentOrFailed()
	if err != nil {
		b.logger.Error(ctx, err)
		return nil, err
	}

	var (
		resp dtos.PageUserWatchedBuddhasResp
		buddhaDOs []*models.Buddha
	)
	queryStream := optionstream.NewQueryStream(req.QueryOptions, req.Limit, req.Offset)
	queryStream.SetOption(queryoptions.EqualUserId, authIdent.GetUserId()).SetOption(queryoptions.IsUnexpired, true)
	buddhaDOs, resp.Pagination, err = services.NewBuddha(b.logger).PageUserWatchedBuddhas(ctx, queryStream)
	if err != nil {
		b.logger.Error(ctx, err)
		return nil, err
	}

	resp.List = transBuddhaDOsToDTOs(buddhaDOs)

    return &resp, nil
}


func (b *Buddha) CreateBuddhaRentOrder(ctx api.Context, req *dtos.CreateBuddhaRentOrderReq) (*dtos.CreateBuddhaRentOrderResp, error) {
    var resp dtos.CreateBuddhaRentOrderResp

    authIdent, err := ctx.GetAuthIdentOrFailed()
	if err != nil {
		b.logger.Error(ctx, err)
		return nil, err
	}

	order, err := services.NewBuddhaRentOrder(b.logger).CreateOrder(ctx, authIdent.GetUserId(), req.RentPackageId, "")
	if err != nil {
		b.logger.Error(ctx, err)
		return nil, err
	}

	resp.Order = transBuddhaRentOrderDOToDTO(order)

    return &resp, nil
}


func (b *Buddha) PageBuddhaRentPackages(ctx api.Context, req *dtos.PageBuddhaRentPackagesReq) (*dtos.PageBuddhaRentPackagesResp, error) {
	authIdent, err := ctx.GetAuthIdentOrFailed()
	if err != nil {
		b.logger.Error(ctx, err)
		return nil, err
	}

	var (
		resp dtos.PageBuddhaRentPackagesResp
		rentPackageDOs []*models.BuddhaRentPackage
	)
	queryStream := optionstream.NewQueryStream(req.QueryOptions, req.Limit, req.Offset)
   	queryStream.SetOption(queryoptions.EqualUserId, authIdent.GetUserId()).
		SetOption(queryoptions.OnShelfStatus, true).
		SetOption(queryoptions.SelectColumns, []string{
			enum.FieldId,
			enum.FieldBuddhaId,
			enum.FieldPrice,
			enum.FieldAvailableDuration,
			enum.FieldName,
			enum.FieldDesc,
			enum.FieldShelfStatus,
		})
	rentPackageDOs, resp.Pagination, err = services.NewBuddha(b.logger).PageBuddhaRentPackages(ctx, queryStream)
	if err != nil {
		b.logger.Error(ctx, err)
		return nil, err
	}

	resp.List = transBuddhaRentPackageDOsToDTOs(rentPackageDOs)

    return &resp, nil
}

func (b *Buddha) PrayToBuddha(ctx api.Context, req *dtos.PrayToBuddhaReq) (*dtos.PrayToBuddhaResp, error) {
    var resp dtos.PrayToBuddhaResp

	authIdent, err := ctx.GetAuthIdentOrFailed()
	if err != nil {
		b.logger.Error(ctx, err)
		return nil, err
	}

	userPrayPropSrv := services.NewUserPrayProp(b.logger)
	var notFreePrayPropIds []uint64
	if len(req.PrayPropIds) > 0 {
		queryStream := optionstream.NewQueryStream(nil, 0, 0)
		queryStream.SetOption(queryoptions.InIds, req.PrayPropIds).
			SetOption(queryoptions.SelectColumns, []string{enum.FieldId, enum.FieldPrice})

		prayPropDOs, _, err := services.NewPrayProp(b.logger).PagePrayProps(ctx, queryStream)
		if err != nil {
			b.logger.Error(ctx, err)
			return nil, err
		}

		for _, prayPropDO := range prayPropDOs {
			if prayPropDO.Price == 0 {
				continue
			}
			notFreePrayPropIds = append(notFreePrayPropIds, prayPropDO.Id)
		}

		if len(notFreePrayPropIds) > 0 {
			ok, err := userPrayPropSrv.EnsurePropsEnough(ctx, authIdent.GetUserId(), notFreePrayPropIds)
			if err != nil {
				b.logger.Error(ctx, err)
				return nil, err
			}

			if !ok {
				b.logger.Error(ctx, err)
				return nil, apperrdef.NewErr(errs.ErrCodeUserPrayPropNotFound)
			}
		}
	}

	_, err = services.NewUserPray(b.logger).CreatePray(ctx, authIdent.GetUserId(), req.BuddhaId, req.Content, req.PrayPropIds)
	if err != nil {
		b.logger.Error(ctx, err)
		return nil, err
	}

	err = userPrayPropSrv.ConsumeProps(ctx, authIdent.GetUserId(), notFreePrayPropIds)
	if err != nil {
		b.logger.Error(ctx, err)
		return nil, err
	}

    return &resp, nil
}

func (b *Buddha) WorshipToBuddha(ctx api.Context, req *dtos.WorshipToBuddhaReq) (*dtos.WorshipToBuddhaResp, error) {
    var resp dtos.WorshipToBuddhaResp

    worshipPropDO, err := services.NewWorshipProp(b.logger).GetWorshipProp(ctx, req.WorshipPropId)
	if err != nil {
		b.logger.Error(ctx, err)
		return nil, err
	}

	authIdent, err := ctx.GetAuthIdentOrFailed()
	if err != nil {
		b.logger.Error(ctx, err)
		return nil, err
	}

	userWorshipSrv := services.NewUserWorshipProp(b.logger)

	if worshipPropDO.IsFree() {
		err = services.NewBuddhaWorship(b.logger).CreateWorship(ctx, &services.CreateWorshipReq{
			UserId: authIdent.GetUserId(),
			BuddhaId: req.BuddhaId,
			WorshipPropDO: worshipPropDO,
		})
		if err != nil {
			b.logger.Error(ctx, err)
			return nil, err
		}

		return &resp, nil
	}

	userWorshipPropDO, err := userWorshipSrv.GetUserWorshipProp(ctx, authIdent.GetUserId(), req.WorshipPropId, req.ConsumeUserWorshipPropId)
	if err != nil {
		b.logger.Error(ctx, err)
		return nil, err
	}

	err = services.NewBuddhaWorship(b.logger).CreateWorship(ctx, &services.CreateWorshipReq{
		UserId: authIdent.GetUserId(),
		BuddhaId: req.BuddhaId,
		WorshipPropDO: worshipPropDO,
		ConsumeUserWorshipPropDO: userWorshipPropDO,
		NeedConsumeUserProp: true,
	})
	if err != nil {
		b.logger.Error(ctx, err)
		return nil, err
	}

	err = userWorshipSrv.ConsumeWorshipProp(ctx, authIdent.GetUserId(), req.ConsumeUserWorshipPropId)
	if err != nil {
		b.logger.Error(ctx, err)
		return nil, err
	}

	return &resp, nil
}