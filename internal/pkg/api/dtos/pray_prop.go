package dtos

type CreatePrayPropOrderReq struct {
	PrayPropId uint64 `json:"pray_prop_id"`
	Num uint32
}

type CreatePrayPropOrderResp struct {
}

type PrayProp struct {
	Name string `json:"name" gorm:"type:varchar"`
	Price uint32 `json:"price"`
	Remark string `json:"remark" gorm:"type:varchar"`
	ShelfStatus int `json:"shelf_status"`
	Extra string `json:"extra"`
}

type PagePrayPropsResp struct {
	PageQueryResp
	List []*PrayProp
}