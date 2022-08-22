package services

import (
	"context"
	"github.com/995933447/apperrdef"
	"github.com/995933447/log-go"
	"github.com/vision-first/wegod/internal/pkg/datamodel/models"
	"github.com/vision-first/wegod/internal/pkg/errs"
	"github.com/vision-first/wegod/internal/pkg/facades"
	"gorm.io/gorm"
)

type UserBuddhaRent struct {
	logger *log.Logger
}

func NewUserBuddhaRent(logger *log.Logger) *UserBuddhaRent {
	return &UserBuddhaRent {
		logger: logger,
	}
}

func (u *UserBuddhaRent) TransBuddhaRentErr(err error) error {
	switch err {
	case gorm.ErrRecordNotFound:
		return apperrdef.NewErr(errs.ErrCodeBuddhaNotRend)
	}
	return err
}

func (u *UserBuddhaRent) GetRendBuddhaExpireAt(ctx context.Context, userId, buddhaId uint64) (int64, error) {
	var rendBuddhaDO models.UserBuddhaRent
	err := facades.MustGORMDB(ctx, u.logger).Where(&models.UserBuddhaRent{
		UserId: userId,
		BuddhaId: buddhaId,
	}).First(&rendBuddhaDO).Error
	if err != nil {
		u.logger.Error(ctx, err)
		return 0, u.TransBuddhaRentErr(err)
	}
	return rendBuddhaDO.ExpireAt, nil
}
