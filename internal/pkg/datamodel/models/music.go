package models

import (
	"github.com/995933447/dbdriverutil/field"
)

type Music struct {
	BaseModel
	AudioUrl string `json:"audio_url" gorm:"type:varchar"`
	BuddhaIds field.Uint64s `json:"buddha_ids"`
	ShelfStatus uint32 `json:"shelf_status"`
}
