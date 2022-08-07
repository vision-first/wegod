package datamodels

import (
	"gorm.io/gorm"
)

//go:generate structfieldconstgen -findPkgPath ../datamodels -outFile ../db/enum/fields.go
type BaseModel struct {
	Id uint64 `gorm:"primaryKey"`
	CreatedAt int64 `gorm:"autoCreateTime"`
	UpdatedAt int64 `gorm:"autoUpdateTime"`
	DeletedAt gorm.DeletedAt `gorm:"index"`
}
