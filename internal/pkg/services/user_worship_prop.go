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
)

type UserWorshipProp struct {
	logger *log.Logger
}

func NewUserWorshipProp(logger *log.Logger) *UserWorshipProp {
	return &UserWorshipProp {
		logger: logger,
	}
}

func (u *UserWorshipProp) TransErr(err error) error {
	switch err {
	case gorm.ErrRecordNotFound:
		return apperrdef.NewErr(errs.ErrCodeUserWorshipPropNotFound)
	}
	return err
}

func (u *UserWorshipProp) PageUserWorshipProps(ctx context.Context, queryStream *optionstream.QueryStream) ([]*models.UserWorshipProp, *optionstream.Pagination, error) {
	db := facades.MustGormDB(ctx, u.logger)

	queryStreamProcessor := optionstream.NewQueryStreamProcessor(queryStream)
	queryStreamProcessor.OnUint64(queryoptions.EqualUserId, func(val uint64) error {
		db.Where(&models.UserWorshipProp{UserId: val})
		return nil
	})

	var userWorshipPropDOs []*models.UserWorshipProp
	pagination, err := queryStreamProcessor.PaginateFrom(ctx, gormimpl.NewOptStreamQuery(db), &userWorshipPropDOs)
	if err != nil {
		u.logger.Error(ctx, err)
		return nil, nil, u.TransErr(err)
	}

	return userWorshipPropDOs, pagination, nil
}

func (u *UserWorshipProp) GetUserWorshipProp(ctx context.Context, userId, worshipPropId, id uint64) (*models.UserWorshipProp, error) {
	var userWorshipPropDO models.UserWorshipProp
	err := facades.MustGormDB(ctx, u.logger).
		Where(map[string]interface{}{enum.FieldId: id, enum.FieldUserId: userId, enum.FieldWorshipPropId: worshipPropId}).
		First(&userWorshipPropDO).
		Error
	if err != nil {
		u.logger.Error(ctx, err)
		return nil, u.TransErr(err)
	}
	return &userWorshipPropDO, nil
}

func (u *UserWorshipProp) ExistUserWorshipProp(ctx context.Context, userId, worshipPropId, id uint64) (bool, error) {
	var userWorshipNum int64
	err := facades.MustGormDB(ctx, u.logger).
		Where(map[string]interface{}{enum.FieldId: id, enum.FieldUserId: userId, enum.FieldWorshipPropId: worshipPropId}).
		Count(&userWorshipNum).
		Error
	if err != nil {
		return false, u.TransErr(err)
	}
	return userWorshipNum > 0, nil
}

func (u *UserWorshipProp) ConsumeUserWorshipProp(ctx context.Context, userId, id uint64) error {
	err := facades.MustGormDB(ctx, u.logger).
		Where(map[string]interface{}{enum.FieldId: id, enum.FieldUserId: userId}).
		Delete(ctx).
		Error
	if err != nil {
		return u.TransErr(err)
	}
	return nil
}