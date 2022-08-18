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

type WorshipProp struct {
	logger *log.Logger
}

func NewWorshipProp(logger *log.Logger) *WorshipProp {
	return &WorshipProp {
		logger: logger,
	}
}

func (*WorshipProp) TransErr(err error) error {
	switch err {
	case gorm.ErrRecordNotFound:
	}
	return err
}

func (w *WorshipProp) PageWorshipProps(ctx context.Context, queryStream *optionstream.QueryStream) ([]*models.WorshipProp, *optionstream.Pagination, error) {
	db := facades.MustGormDB(ctx, w.logger).Order(clause.OrderByColumn{Column: clause.Column{Name: enum.FieldSort}, Desc: true})

	queryStreamProcessor := optionstream.NewQueryStreamProcessor(queryStream)
	queryStreamProcessor.
		OnNone(queryoptions.OnShelfStatus, func() error {
			db.Where(map[string]interface{}{enum.FieldShelfStatus: models.ShelfStatusOnShelf})
			return nil
		}).
		OnUint64List(queryoptions.InIds, func(val []uint64) error {
			db.Where(enum.FieldId + " IN ?", val)
			return nil
		}).
		OnStringList(queryoptions.SelectColumns, gormimpl.MakeOnSelectColumnsOptHandler(db))

	var worshipPropDOs []*models.WorshipProp
	pagination, err := queryStreamProcessor.PaginateFrom(ctx, gormimpl.NewOptStreamQuery(db), &worshipPropDOs)
	if err != nil {
		w.logger.Error(ctx, err)
		return nil, nil, w.TransErr(err)
	}

	return worshipPropDOs, pagination, nil
}

func (w *WorshipProp) IsWorshipPropFree(ctx context.Context, id uint64) (bool, error) {
	var worshipPropDO models.WorshipProp
	err := facades.MustGormDB(ctx, w.logger).Where(map[string]interface{}{
		enum.FieldId: id,
	}).First(&worshipPropDO).Error
	if err != nil {
		w.logger.Error(ctx, err)
		return false, err
	}
	return worshipPropDO.Price == 0, nil
}

func (w *WorshipProp) GetWorshipProp(ctx context.Context, id uint64) (*models.WorshipProp, error)  {
	var prop models.WorshipProp
	facades.MustGormDB(ctx, w.logger).Where(map[string]interface{}{enum.FieldId: id}).First(ctx, &prop)
	return &prop, nil
}

