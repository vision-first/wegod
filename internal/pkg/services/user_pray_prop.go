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

type UserPrayProp struct {
	logger *log.Logger
}

func NewUserPrayProp(logger *log.Logger) *UserPrayProp {
	return &UserPrayProp {
		logger: logger,
	}
}

func (u *UserPrayProp) TransErr(err error) error {
	switch err {
	case gorm.ErrRecordNotFound:
		return apperrdef.NewErr(errs.ErrCodeUserPrayPropNotFound)
	}
	return err
}

func (u *UserPrayProp) EnsureUserEnoughProps(ctx context.Context, userId uint64, prayPropIds []uint64) (bool, error) {
	var propNum int64
	err := facades.MustGormDB(ctx, u.logger).
		Where(&models.UserPrayProp{UserId: userId}).
		Where(enum.FieldPropId + " IN ? AND " + enum.FieldNum + " > 0", prayPropIds).
		Group(enum.FieldPropId).
		Count(&propNum).
		Error
	if err != nil {
		u.logger.Error(ctx, err)
		return false, u.TransErr(err)
	}
	return int(propNum) >= len(prayPropIds), nil
}

func (u *UserPrayProp) ConsumeUserProps(ctx context.Context, userId uint64, prayPropIds []uint64) error {
	err := facades.MustGormDB(ctx, u.logger).
		Where(&models.UserPrayProp{UserId: userId}).
		Where(enum.FieldPropId + " IN ? AND " + enum.FieldNum + " > 0", prayPropIds).
		Update(enum.FieldNum, gorm.Expr(enum.FieldNum + " - 1")).
		Error
	if err != nil {
		u.logger.Error(ctx, err)
		return u.TransErr(err)
	}
	return nil
}

func (u *UserPrayProp) PageUserPrayProps(ctx context.Context, queryStream *optionstream.QueryStream) ([]*models.UserPrayProp, *optionstream.Pagination, error) {
	db := facades.MustGormDB(ctx, u.logger).
		Where(enum.FieldNum + " > 0").
		Order(clause.OrderByColumn{Column: clause.Column{Name: enum.FieldCreatedAt}, Desc: true})

	queryStreamProcessor := optionstream.NewQueryStreamProcessor(queryStream)
	queryStreamProcessor.OnStringList(queryoptions.SelectColumns, gormimpl.MakeOnSelectColumnsOptHandler(db))

	var userPrayPropDOs []*models.UserPrayProp
	pagination, err := queryStreamProcessor.PaginateFrom(ctx, gormimpl.NewOptStreamQuery(db), &userPrayPropDOs)
	if err != nil {
		u.logger.Error(ctx, err)
		return nil, nil, u.TransErr(err)
	}

	return userPrayPropDOs, pagination, nil
}