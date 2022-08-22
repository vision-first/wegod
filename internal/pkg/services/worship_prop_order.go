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

type WorshipPropOrder struct {
	logger *log.Logger
}

func NewWorshipPropOrder(logger *log.Logger) *WorshipPropOrder {
	return &WorshipPropOrder {
		logger: logger,
	}
}

func (w *WorshipPropOrder) TransErr(err error) error {
	switch err {
	case gorm.ErrRecordNotFound:
		return apperrdef.NewErr(errs.ErrCodeWorshipOrderNotFound)
	}
	return err
}

func (w *WorshipPropOrder) CreateOrder(ctx context.Context, userId, worshipPropId uint64) (*models.WorshipPropOrder, error) {
	var worshipPropDO models.WorshipProp
	db := facades.MustGormDB(ctx, w.logger)
	if err := db.First(&worshipPropDO, worshipPropId).Error; err != nil {
		w.logger.Error(ctx, err)
		return nil, w.TransErr(err)
	}
	orderDO := &models.WorshipPropOrder{
		UserId: userId,
		PropId: worshipPropId,
		Sn: bson.NewObjectId().Hex(),
		PropAvailableDurationSnapshot: worshipPropDO.AvailableDuration,
		Status: models.OrderStatusNotPay,
	}
	err := db.Create(orderDO).Error
	if err != nil {
		w.logger.Error(ctx, err)
		return nil, w.TransErr(err)
	}
	return orderDO, nil
}

func (w *WorshipPropOrder) GetOrderBySn(ctx context.Context, userId uint64, sn string) (*models.WorshipPropOrder, error) {
	var orderDO models.WorshipPropOrder
	err := facades.MustGormDB(ctx, w.logger).
		Where(&models.WorshipPropOrder{UserId: userId, Sn: sn}).
		First(ctx, &orderDO).
		Error
	if err != nil {
		w.logger.Error(ctx, err)
		return nil, w.TransErr(err)
	}
	return &orderDO, nil
}



