package models

import (
	"github.com/995933447/dbdriverutil/field"
)

const (
	ProductTypeReal = iota
	ProductTypeVirtual
)

const (
	DeliveryTypePost = iota
	DeliveryTypeFreePost
	DeliveryTypeCustomerService
)

type ShopProduct struct {
	BaseModel
	Name string `json:"name"`
	CategoryId uint64 `json:"category_id"`
	MainImages field.Strings `json:"main_images"`
	Desc string `json:"desc"`
	ProductType uint32 `json:"product_type"`
	DeliveryType uint32 `json:"delivery_type"`
	ShelfStatus uint32 `json:"shelf_status"`
	InventoryNum uint32 `json:"inventory_num"`
	SaleNum uint32 `json:"sale_num" gorm:"index"`
	Price uint32 `json:"price"`
	OnShelfAt int64 `json:"on_shelf_at"`
	OffShelfAt int64 `json:"off_shelf_at"`
}
