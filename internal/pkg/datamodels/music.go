package datamodels

import "github.com/vision-first/wegod/internal/pkg/db"

type Music struct {
	BaseModel
	AudioUrl string `json:"audio_url" gorm:"type:varchar"`
	BuddhaIds db.Uint64s `json:"buddha_ids"`
	ShelfStatus uint32 `json:"shelf_status"`
}
