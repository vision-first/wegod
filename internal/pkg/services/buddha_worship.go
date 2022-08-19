package services

import (
	"context"
	"github.com/995933447/apperrdef"
	"github.com/995933447/log-go"
	"github.com/995933447/optionstream"
	"github.com/vision-first/wegod/internal/pkg/datamodel/models"
	"github.com/vision-first/wegod/internal/pkg/db/mysql/orms/gormimpl"
	"github.com/vision-first/wegod/internal/pkg/errs"
	"github.com/vision-first/wegod/internal/pkg/facades"
	"github.com/vision-first/wegod/internal/pkg/queryoptions"
	"gorm.io/gorm"
	"time"
)

type BuddhaWorship struct {
	logger *log.Logger
}

func NewBuddhaWorship(logger *log.Logger) *BuddhaWorship {
	return &BuddhaWorship {
		logger: logger,
	}
}

func (b *BuddhaWorship) TransErr(err error) error {
	switch err {
	case gorm.ErrRecordNotFound:
	}
	return err
}


func (b *BuddhaWorship) PageBuddhaWorship(ctx context.Context, queryStream *optionstream.QueryStream) ([]*models.BuddhaWorship, *optionstream.Pagination, error) {
	db := facades.MustGormDB(ctx, b.logger)

	queryStreamProcessor := optionstream.NewQueryStreamProcessor(queryStream)
	queryStreamProcessor.
		OnUint64(queryoptions.EqualUserId, func(val uint64) error {
			db.Where(&models.BuddhaWorship{UserId: val})
			return nil
		}).
		OnUint64(queryoptions.EqualBuddhaId, func(val uint64) error {
			db.Where(&models.BuddhaWorship{BuddhaId: val})
			return nil
		})

	var buddhaWorshipDOs []*models.BuddhaWorship
	pagination, err := queryStreamProcessor.PaginateFrom(ctx, gormimpl.NewOptStreamQuery(db), &buddhaWorshipDOs)
	if err != nil {
		b.logger.Error(ctx, err)
		return nil, nil, b.TransErr(err)
	}

	return buddhaWorshipDOs, pagination, nil
}

type CreateWorshipReq struct {
	UserId, BuddhaId uint64
	NeedConsumeUserProp bool
	WorshipPropDO *models.WorshipProp
	ConsumeUserWorshipPropDO *models.UserWorshipProp
}

func (b *BuddhaWorship) CreateWorship(ctx context.Context, req *CreateWorshipReq) error {
	db := facades.MustGormDB(ctx, b.logger)


	if req.NeedConsumeUserProp && req.ConsumeUserWorshipPropDO == nil {
		return apperrdef.NewErrWithMsg(errs.ErrCodeInternal, "req.consumeUserWorshipPropDO is required when req.needConsumeUserProp be true.")
	}

	buddhaWorshipDO := &models.BuddhaWorship{
		BuddhaId: req.BuddhaId,
		UserId: req.UserId,
		WorshipPropId: req.WorshipPropDO.Id,
	}

	if req.NeedConsumeUserProp {
		buddhaWorshipDO.ConsumeUserWorshipId = req.ConsumeUserWorshipPropDO.Id
		buddhaWorshipDO.ExpireAt = req.ConsumeUserWorshipPropDO.AvailableDuration + time.Now().Unix()
	} else {
		buddhaWorshipDO.ExpireAt = req.WorshipPropDO.AvailableDuration + time.Now().Unix()
	}

	if err := db.Create(buddhaWorshipDO).Error; err != nil {
		b.logger.Error(ctx, err)
		return b.TransErr(err)
	}

	return nil
}
