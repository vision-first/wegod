package datamodels

import "gorm.io/gorm"

type Buddha struct {
	Id uint64 `gorm:"primaryKey"`
	CreatedAt int64 `gorm:"autoCreateTime"`
	UpdatedAt int64 `gorm:"autoUpdateTime"`
	DeletedAt gorm.DeletedAt `gorm:"index"`
	Name string
	Image string
	Sort uint32
	Remark string
}