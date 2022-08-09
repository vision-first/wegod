package datamodels

import (
	"github.com/vision-first/wegod/internal/pkg/db"
)

type UserPray struct {
	BaseModel
	UserId uint64 `json:"user_id" gorm:"type:index"`
	BuddhaId uint64 `json:"buddha_id" gorm:"type:index"`
	Content string `json:"content"`
	PrayPropIds db.Uint64s `json:"prop_ids" gorm:"type:json"`
}
