package dtos

import (
	"github.com/995933447/dbdriverutil/field"
)

type ShopProduct struct {
	Name string `json:"name"`
	CategoryId uint64 `json:"category_id"`
	MainImages []string `json:"main_images"`
	Desc string `json:"desc"`
	ProductType uint32 `json:"product_type"`
	DeliveryType uint32 `json:"delivery_type"`
	ShelfStatus uint32 `json:"shelf_status"`
	InventoryNum uint32 `json:"inventory_num"`
	SaleNum uint32 `json:"sale_num" gorm:"index"`
	Price uint32 `json:"price"`
}

type PageShopProductsResp struct {
	PageQueryResp
	List []*ShopProduct
}

type ShopProductCategory struct {
	Name string
	Id uint64
}

type PageShopProductCategoriesResp struct {
	PageQueryResp
	List []*ShopProductCategory
}

type CreateShopOrderReq struct {
	ProductId uint64
	Num uint32
	Remark string
	// 收货人姓名
	ConsigneeName string
	// 收货人的电话号码
	ConsigneePhone string
	// 省份
	ConsigneeProvince string
	// 城市
	ConsigneeCity string
	// 区域
	ConsigneeDistrict string
	// 收货人的详细地址
	ConsigneeAddress string
	// 物流单号
	LogisticsSn string
	// 物流公司
	LogisticsCompany string
}

type CreateShopOrderResp struct {
	Order *ShopOrder
}

type ShopOrder struct {
	Sn string
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
	// 商品数量
	ProductNum uint32 `json:"product_num"`
	// 商品价格
	BoughtNum  uint32 `json:"product_price"`
}

type PageShopOrdersResp struct {
	PageQueryResp
	List []*ShopOrder
}



type GetOrderReq struct {
	OrderId uint64
}

type GetOrderResp struct {
	Order *ShopOrder
}

type GetProductReq struct {
	Id uint64
}

type GetProductResp struct {
	ShopProduct *ShopProduct
}