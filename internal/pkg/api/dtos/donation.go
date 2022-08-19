package dtos

type CreateDonationOrderReq struct {
	BuddhaId uint64 `json:"buddha_id"`
	Money uint32 `json:"money"`
	Scene int `json:"scene"`
}

type CreateDonationOrderResp struct {
}

type PageDonationRanksReq struct {
	PageQueryReq
}

type UserDonationRecord struct {
	NickName string
	Money uint32
}

type PageDonationRanksResp struct {
	PageQueryResp
	List []*UserDonationRecord
}