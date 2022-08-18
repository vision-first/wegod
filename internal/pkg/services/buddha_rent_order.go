package services

import (
	"context"
	"github.com/995933447/apperrdef"
	"github.com/995933447/log-go"
	"github.com/vision-first/wegod/internal/pkg/datamodel/models"
	"github.com/vision-first/wegod/internal/pkg/db/enum"
	"github.com/vision-first/wegod/internal/pkg/errs"
	"github.com/vision-first/wegod/internal/pkg/facades"
	"gopkg.in/mgo.v2/bson"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
	"time"
)

type BuddhaRentOrder struct {
	logger *log.Logger
}

func NewBuddhaRentOrder(logger *log.Logger) *BuddhaRentOrder {
	return &BuddhaRentOrder{
		logger: logger,
	}
}

func (b *BuddhaRentOrder) TransErr(err error) error {
	switch err {
	case gorm.ErrRecordNotFound:
		return apperrdef.NewErr(errs.ErrCodeBuddhaRentOrderNotFound)
	}
	return err
}

func (*BuddhaRentOrder) GenOrderSn() string {
	return bson.NewObjectId().Hex()
}

func (b *BuddhaRentOrder) CreateOrder(ctx context.Context, userId, rentPackageId uint64, remark string) (*models.BuddhaRentOrder, error) {
	db := facades.MustGormDB(ctx, b.logger)

	var rentPackageDO models.BuddhaRentPackage
	if err := db.First(&rentPackageDO, rentPackageId).Error; err != nil {
		b.logger.Error(ctx, err)
		return nil, err
	}

	orderDO := &models.BuddhaRentOrder{
		UserId:                               userId,
		BuddhaId:                             rentPackageDO.BuddhaId,
		RentPackageId:                        rentPackageId,
		RentPackagePriceSnapshot:             rentPackageDO.Price,
		RentPackageAvailableDurationSnapshot: rentPackageDO.AvailableDuration,
		RentPackageNameSnapshot:              rentPackageDO.Name,
		RentPackageDescSnapshot:              rentPackageDO.Desc,
		Status:                               models.OrderStatusNotPay,
		Sn:                                   b.GenOrderSn(),
		Money: 								  rentPackageDO.Price,
		Remark: 							  remark,
	}
	err := db.Create(orderDO).Error
	if err != nil {
		b.logger.Error(ctx, err)
		return nil, err
	}

	return orderDO, nil
}

func (b *BuddhaRentOrder) ExistEffectiveOrder(ctx context.Context, userId, BuddhaId uint64) (bool, error) {
	var unexpiredOrderNum int64
	err := facades.MustGormDB(ctx, b.logger).
		Where(&models.BuddhaRentOrder{UserId: userId, BuddhaId: BuddhaId, Status: models.OrderStatusPayed}).
		Where(enum.FieldEndRentAt + "> ? OR " + enum.FieldEndRentAt + " = 0", time.Now().Unix()).
		Count(&unexpiredOrderNum).
		Error
	if err != nil {
		b.logger.Error(ctx, err)
		return false, b.TransErr(err)
	}
	return unexpiredOrderNum > 0, nil
}

func (b *BuddhaRentOrder) GetLongestAndEffectiveOrder(ctx context.Context) (*models.BuddhaRentOrder, error) {
	var orderDO models.BuddhaRentOrder
	db := facades.MustGormDB(ctx, b.logger)
	err := db.Where(map[string]interface{}{enum.FieldEndRentAt: 0, enum.FieldStatus: models.OrderStatusPayed}).
		First(&orderDO).
		Error
	if err == gorm.ErrRecordNotFound {
		err = db.Where(enum.FieldEndRentAt + " > ?", time.Now().Unix()).
			Where(map[string]interface{}{enum.FieldStatus: models.OrderStatusPayed}).
			Order(clause.OrderByColumn{Column: clause.Column{Name: enum.FieldEndRentAt}, Desc: true}).
			First(&orderDO).
			Error
	}
	if err != nil {
		b.logger.Error(ctx, err)
		return nil, b.TransErr(err)
	}
	return &orderDO, nil
}
