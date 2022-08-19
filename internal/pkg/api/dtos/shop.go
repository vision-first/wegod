package dtos

type PageShopProductsReq struct {
	PageBuddhaReq
}

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

type PageShopProductCategoriesReq struct {
	PageBuddhaReq
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
}

type CreateShopOrderResp struct {
}

type ShopOrder struct {
	Sn string

}

type PageShopOrdersReq struct {
}

type PageShopOrdersResp struct {
}

