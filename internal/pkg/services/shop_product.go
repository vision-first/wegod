package services

import (
	"context"
	"github.com/995933447/apperrdef"
	"github.com/vision-first/wegod/internal/pkg/errs"
	"gorm.io/gorm"
	"github.com/995933447/log-go"
	"github.com/995933447/optionstream"
	"github.com/vision-first/wegod/internal/pkg/datamodel/models"
	"github.com/vision-first/wegod/internal/pkg/db/enum"
	"github.com/vision-first/wegod/internal/pkg/db/mysql/orms/gormimpl"
	"github.com/vision-first/wegod/internal/pkg/facades"
	"github.com/vision-first/wegod/internal/pkg/queryoptions"
	"gorm.io/gorm/clause"
)

type ShopProduct struct {
	logger *log.Logger
}

func NewShopProduct(logger *log.Logger) *ShopProduct {
	return &ShopProduct {
		logger: logger,
	}
}

func (s *ShopProduct) TransErr(err error) error {
	switch err {
	case gorm.ErrRecordNotFound:
		return apperrdef.NewErr(errs.ErrCodeShopProductNotFound)
	}
	return err
}

func (s *ShopProduct) GetProductById(ctx context.Context, id uint64) (*models.ShopProduct, error) {
	var product models.ShopProduct
	if err := facades.MustGORMDB(ctx, s.logger).First(&product, id).Error; err != nil {
		s.logger.Error(ctx, err)
		return nil, s.TransErr(err)
	}
	return &product, nil
}

func (s *ShopProduct) PageProducts(ctx context.Context, queryStream *optionstream.QueryStream) ([]*models.ShopProduct, *optionstream.Pagination, error) {
	db := facades.MustGORMDB(ctx, s.logger)

	queryStreamProcessor := optionstream.NewQueryStreamProcessor(queryStream)
	queryStreamProcessor.
		OnString(queryoptions.LikeName, func(val string) error {
			db.Where(enum.FieldName + " Like %?%", val)
			return nil
		}).
		OnNone(queryoptions.OnShelfStatus, func() error {
			db.Where(&models.ShopProduct{ShelfStatus: models.ShelfStatusOnShelf})
			return nil
		}).
		OnUint64(queryoptions.EqualCategoryId, func(val uint64) error {
			db.Where(&models.ShopProduct{CategoryId: val})
			return nil
		}).
		OnNone(queryoptions.OrderByShelfAtDesc, func() error {
			db.Order(clause.OrderByColumn{
				Column: clause.Column{
					Name: enum.FieldOnShelfAt,
				},
				Desc: true,
			})
			return nil
		}).
		OnNone(queryoptions.OrderBySaleNumDesc, func() error {
			db.Order(clause.OrderByColumn{
				Column: clause.Column{
					Name: enum.FieldSaleNum,
				},
				Desc: true,
			})
			return nil
		})

	var shopProductDOs []*models.ShopProduct
	pagination, err := queryStreamProcessor.PaginateFrom(ctx, gormimpl.NewOptStreamQuery(db), &shopProductDOs)
	if err != nil {
		s.logger.Error(ctx, err)
		return nil, nil, s.TransErr(err)
	}

	return shopProductDOs, pagination, nil
}