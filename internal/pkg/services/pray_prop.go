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
)


type PrayProp struct {
	logger *log.Logger
}

func NewPrayProp(logger *log.Logger) *PrayProp {
	return &PrayProp{
		logger: logger,
	}
}

func (p *PrayProp) TransErr(err error) error {
	switch err {
	case gorm.ErrRecordNotFound:
		return apperrdef.NewErr(errs.ErrCodePrayPropNotFound)
	}
	return err
}

func (p *PrayProp) PagePrayProps(ctx context.Context, queryStream *optionstream.QueryStream) ([]*models.PrayProp, *optionstream.Pagination, error) {
	db := facades.MustGormDB(ctx, p.logger).Order(clause.OrderByColumn{Column: clause.Column{Name: enum.FieldSort}, Desc: true})

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

	var prayPropDOs []*models.PrayProp
	pagination, err := queryStreamProcessor.PaginateFrom(ctx, gormimpl.NewOptStreamQuery(db), &prayPropDOs)
	if err != nil {
		p.logger.Error(ctx, err)
		return nil, nil, p.TransErr(err)
	}

	return prayPropDOs, pagination, nil
}

func (p *PrayProp) IsPrayPropFree(ctx context.Context, id uint64) (bool, error) {
	var prayPropDO models.PrayProp
	err := facades.MustGormDB(ctx, p.logger).
		Select(enum.FieldPrice).
		Where(map[string]interface{}{enum.FieldId: id}).
		First(&prayPropDO).
		Error
	if err != nil {
		p.logger.Error(ctx, err)
		return false, p.TransErr(err)
	}
	return prayPropDO.Price == 0, nil
}