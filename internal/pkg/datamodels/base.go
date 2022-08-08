package datamodels

import (
	"gorm.io/plugin/soft_delete"
)

//go:generate structfieldconstgen -findPkgPath ../datamodels -outFile ../db/enum/fields.go -prefix Field
type BaseModel struct {
	Id uint64 `gorm:"primaryKey"`
	CreatedAt int64 `gorm:"autoCreateTime"`
	UpdatedAt int64 `gorm:"autoUpdateTime"`
	DeletedAt soft_delete.DeletedAt `gorm:"index"`
}
