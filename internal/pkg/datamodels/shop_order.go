package datamodels

import (
	"github.com/995933447/dbdriverutil/field"
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

type ShopOrder struct {
	BaseModel
	OrderSn string `json:"order_sn" gorm:"type:varchar, UniqueIndex"`
	UserId uint64 `json:"user_id" gorm:"index"`
	ProductId uint64 `json:"product_id" gorm:"index:index"`
	ProductNameSnapshot string `json:"product_name_snapshot"`
	ProductMainImagesSnapshot field.Strings `json:"product_main_images_snapshot"`
	Price uint32 `json:"price"`
	DeliveryTypeSnapshot uint32 `json:"delivery_type_snapshot"`
	ProductTypeSnapshot uint32 `json:"product_type_snapshot"`
	Status uint32 `json:"status"`
}
