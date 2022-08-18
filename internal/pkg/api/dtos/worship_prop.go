package dtos

type WorshipProp struct {
	Id uint64 `json:"id"`
	Name string `json:"name" gorm:"type:varchar"`
	Image string `json:"image" gorm:"type:varchar"`
	Price uint32 `json:"price"`
	ShelfStatus       uint `json:"shelf_status"`
	AvailableDuration int64 `json:"available_duration"`
}

type PageWorshipPropsReq struct {
	PageQueryReq
}

type PageWorshipPropsResp struct {
	PageQueryResp
	List []*WorshipProp
}

type CreateWorshipPropOrderReq struct {
	WorshipPropId uint64
	Num uint32
}

type CreateWorshipPropOrderResp struct {
}