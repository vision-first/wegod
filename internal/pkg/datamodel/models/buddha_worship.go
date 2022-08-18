package models

type BuddhaWorship struct {
	BaseModel
	BuddhaId uint64 `json:"buddha_id" gorm:"uniqueIndex"`
	UserId uint64 `json:"user_id" gorm:"uniqueIndex"`
	WorshipPropId uint64 `json:"prop_id"`
	ConsumeUserWorshipId uint64 `json:"consume_user_worship_id"`
	ExpireAt int64 `json:"expire_at"`
}
