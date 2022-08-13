package datamodels

import (
	"github.com/995933447/dbdriverutil/field"
)

type UserPray struct {
	BaseModel
	UserId uint64 `json:"user_id" gorm:"type:index"`
	BuddhaId uint64 `json:"buddha_id" gorm:"type:index"`
	Content string `json:"content"`
	PrayPropIds field.Uint64s `json:"prop_ids" gorm:"type:json"`
}
