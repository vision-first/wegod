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
	"gopkg.in/mgo.v2/bson"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type DonationOrder struct {
	logger *log.Logger
}

func NewDonationOrder(logger *log.Logger) *DonationOrder {
	return &DonationOrder {
		logger: logger,
	}
}

func (d *DonationOrder) TransErr(err error) error {
	switch err {
	case gorm.ErrRecordNotFound:
	}
	return err
}

func (d *DonationOrder) PageOrders(ctx context.Context, queryStream *optionstream.QueryStream) ([]*models.DonationOrder, *optionstream.Pagination, error) {
	db := facades.MustGORMDB(ctx, d.logger).Order(clause.OrderByColumn{Column: clause.Column{Name: enum.FieldId}, Desc: true})

	queryStreamProcessor := optionstream.NewQueryStreamProcessor(queryStream)
	queryStreamProcessor.OnUint64(queryoptions.EqualUserId, func(val uint64) error {
		db.Where(&models.DonationOrder{UserId: val})
		return nil
	})

	var donationOrderDOs []*models.DonationOrder
	pagination, err := queryStreamProcessor.PaginateFrom(ctx, gormimpl.NewOptStreamQuery(db), &donationOrderDOs)
	if err != nil {
		d.logger.Error(ctx, err)
		return nil, nil, d.TransErr(err)
	}

	return donationOrderDOs, pagination, nil
}

func (d *DonationOrder) GetOrderBySn(ctx context.Context, sn string) (*models.DonationOrder, error) {
	var order models.DonationOrder
	err := facades.GORMDB(ctx, d.logger).Where(map[string]interface{}{enum.FieldSn: sn}).First(&order).Error
	if err != nil {
		d.logger.Error(ctx, err)
		return nil, d.TransErr(err)
	}
	return &order, nil
}

type CreateOrderReq struct {
	UserId, BuddhaId uint64
	Money uint32
	Scene models.DonationScene
}

func (d *DonationOrder) CreateOrder(ctx context.Context, req *CreateOrderReq) (*models.DonationOrder, error) {
	order := &models.DonationOrder{
		UserId: req.UserId,
		BuddhaId: req.BuddhaId,
		Money: req.Money,
		Sn: bson.NewObjectId().Hex(),
		DonationScene: uint32(req.Scene),
		Status: models.OrderStatusNotPay,
	}
	if err := facades.MustGORMDB(ctx, d.logger).Create(order).Error; err != nil {
		d.logger.Error(ctx, err)
		return nil, d.TransErr(err)
	}
	return order, nil
}