package services

import (
	"context"
	"github.com/995933447/apperrdef"
	"github.com/995933447/log-go"
	"github.com/995933447/optionstream"
	"github.com/vision-first/wegod/internal/pkg/datamodel/models"
	"github.com/vision-first/wegod/internal/pkg/db/enum"
	"github.com/vision-first/wegod/internal/pkg/db/mysql/orms/gormimpl"
	"github.com/vision-first/wegod/internal/pkg/errs"
	"github.com/vision-first/wegod/internal/pkg/facades"
	"github.com/vision-first/wegod/internal/pkg/queryoptions"
	"gopkg.in/mgo.v2/bson"
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
		return apperrdef.NewErr(errs.ErrCodeShopProductNotFound)
	}
	return err
}

func (s *ShopOrder) GetOrderById(ctx context.Context, id uint64) (*models.ShopOrder, error) {
	var orderDO models.ShopOrder
	if err := facades.MustGORMDB(ctx, s.logger).First(&orderDO, id).Error; err != nil {
		s.logger.Error(ctx, err)
		return nil, s.TransErr(err)
	}
	return &orderDO, nil
}

func (s *ShopOrder) GetOrder(ctx context.Context, optionStream *optionstream.Stream) (*models.ShopOrder, error) {
	db := facades.MustGORMDB(ctx, s.logger)

	err := optionstream.NewStreamProcessor(optionStream).
		OnUint64(queryoptions.EqualUserId, func(val uint64) error {
			db.Where(&models.ShopOrder{UserId: val})
			return nil
		}).
		OnUint64(queryoptions.EqualId, func(val uint64) error {
			db.Where(map[string]interface{}{enum.FieldId: val})
			return nil
		}).
		Process()
	if err != nil {
		s.logger.Error(ctx, err)
		return nil, s.TransErr(err)
	}

	var orderDO models.ShopOrder
	if err = db.First(&orderDO).Error; err != nil {
		s.logger.Error(ctx, err)
		return nil, s.TransErr(err)
	}

	return &orderDO, nil
}

func (s *ShopOrder) PageOrders(ctx context.Context, queryStream *optionstream.QueryStream) ([]*models.ShopOrder, *optionstream.Pagination, error) {
	db := facades.MustGORMDB(ctx, s.logger).Order(clause.OrderByColumn{Column: clause.Column{Name: enum.FieldCreatedAt}, Desc: true})

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

type Consignee struct {
	// 收货人姓名
	ConsigneeName string
	// 收货人的电话号码
	ConsigneePhone string
	// 省份
	ConsigneeProvince string
	// 城市
	ConsigneeCity string
	// 区域
	ConsigneeDistrict string
	// 收货人的详细地址
	ConsigneeAddress string
	// 物流单号
	LogisticsSn string
	// 物流公司
	LogisticsCompany string
}

func (s *ShopOrder) CreateOrder(ctx context.Context, userId, productId uint64, boughtNum uint32, remark string, consignee Consignee) (*models.ShopOrder, error) {
	var productDO models.ShopProduct
	db := facades.MustGORMDB(ctx, s.logger)
	if err := db.First(&productDO, productId).Error; err != nil {
		s.logger.Error(ctx, err)
		return nil, s.TransErr(err)
	}
	orderDO := &models.ShopOrder{
		UserId:        userId,
		Sn:       	   bson.NewObjectId().Hex(),
		ProductPriceSnapshot: productDO.Price,
		Status:        models.OrderStatusPayed,
		Money:         productDO.Price * boughtNum,
		BoughtNum: 	   boughtNum,
		ProductMainImagesSnapshot: productDO.MainImages,
		ProductNameSnapshot: productDO.Name,
		DeliveryTypeSnapshot: productDO.DeliveryType,
		ProductTypeSnapshot: productDO.ProductType,
	}
	err := db.Create(orderDO).Error
	if err != nil {
		s.logger.Error(ctx, err)
		return nil, s.TransErr(err)
	}
	return orderDO, nil
}