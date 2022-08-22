package services

import (
	"context"
	"github.com/995933447/apperrdef"
	"github.com/995933447/log-go"
	"github.com/vision-first/wegod/internal/pkg/datamodel/models"
	"github.com/vision-first/wegod/internal/pkg/errs"
	"github.com/vision-first/wegod/internal/pkg/facades"
	"gopkg.in/mgo.v2/bson"
	"gorm.io/gorm"
)


type PrayPropOrder struct {
	logger *log.Logger
}

func NewPrayPropOrder(logger *log.Logger) *PrayPropOrder {
	return &PrayPropOrder{
		logger: logger,
	}
}

func (p *PrayPropOrder) TransErr(err error) error {
	switch err {
	case gorm.ErrRecordNotFound:
		return apperrdef.NewErr(errs.ErrCodePrayPropOrderNotFound)
	}
	return err
}

func (p *PrayPropOrder) CreateOrder(ctx context.Context, userId, prayPropId uint64, num uint32) (*models.PrayPropOrder, error) {
	var prayPropDO models.PrayProp
	db := facades.MustGORMDB(ctx, p.logger)
	if err := db.First(&prayPropDO, prayPropId).Error; err != nil {
		p.logger.Error(ctx, err)
		return nil, p.TransErr(err)
	}
	orderDO := &models.PrayPropOrder{
		UserId:        userId,
		PropId:        prayPropId,
		PropPriceSnapshot: prayPropDO.Price,
		Status:        models.OrderStatusPayed,
		Sn:            bson.NewObjectId().Hex(),
		Money:         prayPropDO.Price * num,
		Num: 		   num,
	}
	err := db.Create(orderDO).Error
	if err != nil {
		p.logger.Error(ctx, err)
		return nil, p.TransErr(err)
	}
	return orderDO, nil
}

func (p *PrayPropOrder) GetOrderBySn(ctx context.Context, userId uint64, sn string) (*models.PrayPropOrder, error) {
	var orderDO models.PrayPropOrder
	err := facades.MustGORMDB(ctx, p.logger).
		Where(&models.PrayPropOrder{UserId: userId, Sn: sn}).
		First(ctx, &orderDO).
		Error
	if err != nil {
		p.logger.Error(ctx, err)
		return nil, p.TransErr(err)
	}
	return &orderDO, nil
}
