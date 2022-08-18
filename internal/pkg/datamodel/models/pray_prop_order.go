package models

type PrayPropOrder struct {
	BaseModel
	UserId uint64 `json:"user_id"`
	PropId uint64 `json:"prop_id"`
	PayedAt int64 `json:"payed_at"`
	PropPriceSnapshot uint32 `json:"price_snapshot"`
	Status int
	Num uint32
	Money uint32
	Sn string
}
