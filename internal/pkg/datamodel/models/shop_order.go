package models

import (
	"github.com/995933447/dbdriverutil/field"
)

type ShopOrder struct {
	BaseModel
	Sn string `json:"order_sn" gorm:"type:varchar, UniqueIndex"`
	UserId uint64 `json:"user_id" gorm:"index"`
	ProductId uint64 `json:"product_id" gorm:"index:index"`
	ProductNameSnapshot string `json:"product_name_snapshot"`
	ProductMainImagesSnapshot field.Strings `json:"product_main_images_snapshot"`
	Money uint32 `json:"price"`
	DeliveryTypeSnapshot uint32 `json:"delivery_type_snapshot"`
	ProductTypeSnapshot uint32 `json:"product_type_snapshot"`
	Status uint `json:"status"`
	ProductPriceSnapshot uint32 `json:"product_price_snapshot"`
	PayedAt int64 `json:"payed_at"`
	// 备注
	Remark string `json:"remark"`
	// 收货人姓名
	ConsigneeName string `json:"consignee_name"`
	// 收货人的电话号码
	ConsigneePhone string `json:"consignee_phone"`
	// 省份
	ConsigneeProvince string `json:"consignee_province"`
	// 城市
	ConsigneeCity string `json:"consignee_city"`
	// 区域
	ConsigneeDistrict string `json:"consignee_district"`
	// 收货人的详细地址
	ConsigneeAddress string `json:"consignee_address"`
	// 发货时间
	DeliveredAt uint32 `json:"delivered_at"`
	// 完成时间
	FinishAt uint32 `json:"finish_at"`
	// 物流单号
	LogisticsSn string `json:"logistics_sn"`
	// 物流公司
	LogisticsCompany string `json:"logistics_company"`
	// 购买数量
	BoughtNum uint32 `json:"bought_num"`
}
