package services

import (
	"context"
	"github.com/995933447/log-go"
	"github.com/995933447/optionstream"
	"github.com/vision-first/wegod/internal/pkg/datamodels"
	"github.com/vision-first/wegod/internal/pkg/db/enum"
	"github.com/vision-first/wegod/internal/pkg/db/mysql/orms/gormimpl"
	"github.com/vision-first/wegod/internal/pkg/facades"
	"gorm.io/gorm/clause"
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

func (b *Buddha) PageBuddha(ctx context.Context, queryStream *optionstream.QueryStream) ([]*datamodels.Buddha, *optionstream.Pagination, error) {
	var list []*datamodels.Buddha
	db := facades.MustGormDB(ctx, b.logger).Order(clause.OrderByColumn{Column: clause.Column{Name: enum.FieldId}, Desc: true})
	pagination, err := optionstream.NewQueryStreamProcessor(queryStream).PaginateFrom(ctx, gormimpl.NewOptStreamQuery(db), &list)
	if err != nil {
		b.logger.Error(ctx, err)
		return nil, nil, err
	}
	return list, pagination, nil
}

func (b *Buddha) WatchBuddha(ctx context.Context, userId, buddhaId uint64) error {
	buddhaFollow := datamodels.BuddhaFollow {
		UserId: userId,
		BuddhaId: buddhaId,
		WatchedAt: time.Now().Unix(),
	}
	res := facades.MustGormDB(ctx, b.logger).
		Unscoped().
		Where(&datamodels.BuddhaFollow{UserId: userId, BuddhaId: buddhaId}).
		FirstOrCreate(ctx, &buddhaFollow)
	if res.Error != nil {
		return res.Error
	}

	// created
	if res.RowsAffected > 0 {
		return nil
	}

	err := facades.MustGormDB(ctx, b.logger).
		Unscoped().
		Where(map[string]interface{}{enum.FieldId: buddhaFollow.Id}).
		Updates(map[string]interface{}{enum.FieldWatchedAt: time.Now().Unix(), enum.FieldDeletedAt: 0}).
		Error
	if err != nil {
		return err
	}

	return nil
}

func (b *Buddha) UnwatchBuddha(ctx context.Context, userId, buddhaId uint64) error {
	err := facades.MustGormDB(ctx, b.logger).
		Where(&datamodels.BuddhaFollow{UserId: userId, BuddhaId: buddhaId}).
		Delete(&datamodels.BuddhaFollow{}).
		Error
	if err != nil {
		return err
	}
	return nil
}

func (b *Buddha) PageUserWatchedBuddha(ctx context.Context, queryStream *optionstream.QueryStream) ([]*datamodels.Buddha, error) {
	return nil, nil
}