package dtos

type GetUserInfoReq struct {
}

type GetUserInfoResp struct {
	NickName string `json:"nick_name" gorm:"type:varchar"`
	Avatar string `json:"avatar" gorm:"type:varchar"`
	Desc string `json:"desc"`
	Gender uint8 `json:"gender"`
}

type SetUserInfoReq struct {
	NickName string `json:"nick_name" gorm:"type:varchar"`
	Avatar string `json:"avatar" gorm:"type:varchar"`
	Desc string `json:"desc"`
	Gender uint8 `json:"gender"`
}

type SetUserInfoResp struct {
}