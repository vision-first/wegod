package services

import (
	"context"
	"github.com/995933447/log-go"
	"github.com/995933447/optionstream"
	"github.com/vision-first/wegod/internal/pkg/datamodel/models"
	"github.com/vision-first/wegod/internal/pkg/db/enum"
	"github.com/vision-first/wegod/internal/pkg/db/mysql/orms/gormimpl"
	"github.com/vision-first/wegod/internal/pkg/facades"
	"github.com/vision-first/wegod/internal/pkg/queryoptions"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type ShopOrder struct {
	logger *log.Logger
}

func NewShopOrder(logger *log.Logger) *ShopOrder {
	return &ShopOrder {
		logger: logger,
	}
}

func (s *ShopOrder) TransErr(err error) error {
	switch err {
	case gorm.ErrRecordNotFound:
	}
	return err
}

func (s *ShopOrder) PageOrders(ctx context.Context, queryStream *optionstream.QueryStream) ([]*models.ShopOrder, *optionstream.Pagination, error) {
	db := facades.MustGormDB(ctx, s.logger).Order(clause.OrderByColumn{Column: clause.Column{Name: enum.FieldCreatedAt}, Desc: true})

	queryStreamProcessor := optionstream.NewQueryStreamProcessor(queryStream)
	queryStreamProcessor.
		OnUint64(queryoptions.EqualUserId, func(val uint64) error {
			db.Where(&models.ShopOrder{UserId: val})
			return nil
		}).
		OnUint32(queryoptions.EqualStatus, func(val uint32) error {
			db.Where(&models.ShopOrder{Status: uint(val)})
			return nil
		}).
		OnTimestampRange(queryoptions.CreatedAtRange, func(beginAt, endAt int64) error {
			db.Where(enum.FieldCreatedAt + " >= ? and " + enum.FieldCreatedAt + " <= ?", beginAt, endAt)
			return nil
		}).
		OnTimestampRange(queryoptions.PayedAtRange, func(beginAt, endAt int64) error {
			db.Where(enum.FieldPayedAt + " >= ? and " + enum.FieldPayedAt + " <= ?", beginAt, endAt)
			return nil
		})

	var shopOrderDOs []*models.ShopOrder
	pagination, err := queryStreamProcessor.PaginateFrom(ctx, gormimpl.NewOptStreamQuery(db), &shopOrderDOs)
	if err != nil {
		s.logger.Error(ctx, err)
		return nil, nil, s.TransErr(err)
	}

	return shopOrderDOs, pagination, nil
}