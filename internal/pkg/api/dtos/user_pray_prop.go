package dtos

type UserPrayProp struct {
	PrayProp
	Num uint32
}

type PageUserPrayPropsResp struct {
	PageQueryResp
	List []*UserPrayProp
}

