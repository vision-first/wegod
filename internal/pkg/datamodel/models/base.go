package models

import (
	"gorm.io/plugin/soft_delete"
)

const (
	ShelfStatusOnShelf = iota
	ShelfStatusOffShelf
)

const (
	OrderStatusNotPay = iota
	OrderStatusPayed
	OrderStatusDelivered
	OrderStatusFinish
	OrderStatusApplyRefund
	OrderStatusRefunding
	OrderStatusRefunded
)

//go:generate structfieldconstgen -findPkgPath ../models -outFile ../metadata/field_consts.go -prefix Field -func camel
//go:generate structfieldconstgen -findPkgPath ../models -outFile ../../db/enum/fields.go -prefix Field
type BaseModel struct {
	Id uint64 `gorm:"primaryKey"`
	CreatedAt int64 `gorm:"autoCreateTime"`
	UpdatedAt int64 `gorm:"autoUpdateTime"`
	DeletedAt soft_delete.DeletedAt `gorm:"index"`
}
