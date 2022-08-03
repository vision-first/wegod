package services

import (
	"context"
	"github.com/995933447/log-go"
	"github.com/995933447/optionstream"
	"github.com/vision-first/wegod/internal/pkg/datamodels"
	"github.com/vision-first/wegod/internal/pkg/db/mysql/orms/gormimpl"
	"github.com/vision-first/wegod/internal/pkg/facades"
	"gorm.io/gorm/clause"
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
	db := facades.MustGormDB(ctx, b.logger).Order(clause.OrderByColumn{Column: clause.Column{Name: "id"}, Desc: true})
	pagination, err := optionstream.NewQueryStreamProcessor(queryStream).PaginateFrom(ctx, gormimpl.NewOptStreamQuery(db), &list)
	if err != nil {
		b.logger.Error(ctx, err)
		return nil, nil, err
	}
	return list, pagination, nil
}

func (b *Buddha) PleaseBuddha(ctx context.Context, user datamodels.User, buddhaId uint64) error {
	return nil
}