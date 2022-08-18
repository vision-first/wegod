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
	"gorm.io/gorm"
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

func (b *Buddha) TransErr(err error) error {
	switch err {
	case gorm.ErrRecordNotFound:
		err = apperrdef.NewErr(errs.ErrCodeBuddhaNotFound)
	}
	return err
}

func (b *Buddha) PageBuddha(ctx context.Context, queryStream *optionstream.QueryStream) ([]*models.Buddha, *optionstream.Pagination, error) {
	var BuddhaDOs []*models.Buddha

	db := facades.MustGormDB(ctx, b.logger).Order(clause.OrderByColumn{Column: clause.Column{Name: enum.FieldId}, Desc: true})

	optProcessor := optionstream.NewQueryStreamProcessor(queryStream)
	optProcessor.OnStringList(queryoptions.SelectColumns, gormimpl.MakeOnSelectColumnsOptHandler(db))
	pagination, err := optProcessor.PaginateFrom(ctx, gormimpl.NewOptStreamQuery(db), &BuddhaDOs)
	if err != nil {
		b.logger.Error(ctx, err)
		return nil, nil, b.TransErr(err)
	}

	return BuddhaDOs, pagination, nil
}

func (b *Buddha) IsBuddhaFreeRent(ctx context.Context, id uint64) (bool, error) {
	isFreeRentRes := &struct {IsFreeRent bool `json:"is_free_rent"`}{}
	err := facades.MustGormDB(ctx, b.logger).
		Where(map[string]interface{}{enum.FieldId: id}).Select(enum.FieldIsFreeRent).
		First(isFreeRentRes).
		Error
	if err != nil {
		if err != gorm.ErrRecordNotFound {
			b.logger.Error(ctx, err)
			return false, b.TransErr(err)
		}

		isFreeRentRes.IsFreeRent = true
	}

	return isFreeRentRes.IsFreeRent, nil
}

func (b *Buddha) WatchBuddha(ctx context.Context, userId, buddhaId uint64, expireAt int64) error {
	db := facades.MustGormDB(ctx, b.logger)

	buddhaFollowDO := models.BuddhaFollow {
		UserId: userId,
		BuddhaId: buddhaId,
		WatchedAt: time.Now().Unix(),
		ExpireAt: expireAt,
	}
	res := db.Unscoped().Where(&models.BuddhaFollow{UserId: userId, BuddhaId: buddhaId}).FirstOrCreate(ctx, &buddhaFollowDO)
	if res.Error != nil {
		return b.TransErr(res.Error)
	}

	// created
	if res.RowsAffected > 0 {
		return nil
	}

	err := db.Unscoped().
		Where(map[string]interface{}{enum.FieldId: buddhaFollowDO.Id}).
		Updates(map[string]interface{}{
			enum.FieldWatchedAt: time.Now().Unix(),
			enum.FieldDeletedAt: 0,
			enum.FieldExpireAt: expireAt,
		}).
		Error
	if err != nil {
		return b.TransErr(err)
	}

	return nil
}

func (b *Buddha) UnwatchBuddha(ctx context.Context, userId, buddhaId uint64) error {
	err := facades.MustGormDB(ctx, b.logger).
		Where(&models.BuddhaFollow{UserId: userId, BuddhaId: buddhaId}).
		Delete(&models.BuddhaFollow{}).
		Error
	if err != nil {
		b.logger.Error(ctx, err)
		return b.TransErr(err)
	}
	return nil
}

func (b *Buddha) PageUserWatchedBuddhas(ctx context.Context, queryStream *optionstream.QueryStream) ([]*models.Buddha, *optionstream.Pagination, error) {
	followTableName := (&models.BuddhaFollow{}).TableName()
	buddhaTableName := (&models.Buddha{}).TableName()
	db := facades.MustGormDB(ctx, b.logger).
		Model(&models.BuddhaFollow{}).
		Select(buddhaTableName + ".*").
		Joins("RIGHT JOIN " + buddhaTableName +
			" ON " + buddhaTableName + "." + enum.FieldId + " = " + followTableName + "." + enum.FieldBuddhaId).
		Order(followTableName + "." + enum.FieldCreatedAt + " DESC")

	queryStreamProcessor := optionstream.NewQueryStreamProcessor(queryStream)
	queryStreamProcessor.
		OnUint64(queryoptions.EqualUserId, func(val uint64) error {
			db.Where(map[string]interface{}{
				enum.FieldUserId: val,
			})
			return nil
		}).
		OnBool(queryoptions.IsUnexpired, func(val bool) error {
			db.Where(enum.FieldExpireAt + " = 0 OR " + enum.FieldExpireAt + " >= ?", time.Now().Unix())
			return nil
		})

	var buddhas []*models.Buddha
	pagination, err := queryStreamProcessor.PaginateFrom(ctx, gormimpl.NewOptStreamQuery(db), &buddhas)
	if err != nil {
		return nil, nil, err
	}

	return buddhas, pagination, nil
}

func (b *Buddha) PageBuddhaRentPackages(ctx context.Context, queryStream *optionstream.QueryStream) ([]*models.BuddhaRentPackage, *optionstream.Pagination, error) {
	db := facades.MustGormDB(ctx, b.logger).Order(clause.OrderByColumn{Column: clause.Column{Name: enum.FieldSort}, Desc: true})
	queryStreamProcessor := optionstream.NewQueryStreamProcessor(queryStream)
	queryStreamProcessor.
		OnUint64(queryoptions.EqualBuddhaId, func(val uint64) error {
			db.Where(&models.BuddhaRentPackage{BuddhaId: val})
			return nil
		}).
		OnStringList(queryoptions.SelectColumns, gormimpl.MakeOnSelectColumnsOptHandler(db)).
		OnNone(queryoptions.OnShelfStatus, func() error {
			db.Where(&models.BuddhaRentPackage{ShelfStatus: models.ShelfStatusOnShelf})
			return nil
		})

	var rentPackageDOs []*models.BuddhaRentPackage
	pagination, err := queryStreamProcessor.PaginateFrom(ctx, gormimpl.NewOptStreamQuery(db), &rentPackageDOs)
	if err != nil {
		b.logger.Error(ctx, err)
		return nil, nil, err
	}
	return rentPackageDOs, pagination, nil
}

func (b *Buddha) PrayToBuddha(ctx context.Context, userId, buddhaId uint64, prayPropIds []uint64, content string) error {
	err := facades.MustGormDB(ctx, b.logger).Create(&models.UserPray{
		UserId: userId,
		BuddhaId: buddhaId,
		Content: content,
		PrayPropIds: prayPropIds,
	}).Error
	if err != nil {
		b.logger.Error(ctx, err)
		return err
	}
	return nil
}