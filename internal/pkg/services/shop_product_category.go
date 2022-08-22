package services

import (
	"context"
	"github.com/995933447/log-go"
	"github.com/995933447/optionstream"
	"github.com/vision-first/wegod/internal/pkg/datamodel/models"
	"github.com/vision-first/wegod/internal/pkg/db/mysql/orms/gormimpl"
	"github.com/vision-first/wegod/internal/pkg/facades"
	"gorm.io/gorm"
)

type ShopProductCategory struct {
	logger *log.Logger
}

func NewShopProductCategory(logger *log.Logger) *ShopProductCategory {
	return &ShopProductCategory {
		logger: logger,
	}
}

func (s *ShopProductCategory) TransErr(err error) error {
	switch err {
	case gorm.ErrRecordNotFound:
	}
	return err
}

func (s *ShopProductCategory) PageProductCategories(ctx context.Context, queryStream *optionstream.QueryStream) ([]*models.ShopProductCategory, *optionstream.Pagination, error) {
	db := facades.MustGORMDB(ctx, s.logger)

	queryStreamProcessor := optionstream.NewQueryStreamProcessor(queryStream)
	// TODO. set option handler
	var shopProductCategoryDOs []*models.ShopProductCategory
	pagination, err := queryStreamProcessor.PaginateFrom(ctx, gormimpl.NewOptStreamQuery(db), &shopProductCategoryDOs)
	if err != nil {
		s.logger.Error(ctx, err)
		return nil, nil, s.TransErr(err)
	}

	return shopProductCategoryDOs, pagination, nil
}