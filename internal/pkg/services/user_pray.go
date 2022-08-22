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

type UserPray struct {
	logger *log.Logger
}

func NewUserPray(logger *log.Logger) *UserPray {
	return &UserPray {
		logger: logger,
	}
}

func (u *UserPray) TransErr(err error) error {
	switch err {
	case gorm.ErrRecordNotFound:
		return apperrdef.NewErr(errs.ErrCodeUserPrayNotFound)
	}
	return err
}

func (u *UserPray) CreateUserPray(ctx context.Context, userId, buddhaId uint64, content string, prayPropIds []uint64) (*models.UserPray, error) {
	userPray := &models.UserPray{
		UserId: userId,
		Content: content,
		BuddhaId: buddhaId,
		PrayPropIds: prayPropIds,
	}
	if err := facades.MustGORMDB(ctx, u.logger).Create(userPray).Error; err != nil {
		u.logger.Error(ctx, err)
		return nil, u.TransErr(err)
	}
	return userPray, nil
}


