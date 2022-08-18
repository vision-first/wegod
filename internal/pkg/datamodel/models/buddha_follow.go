package models

type BuddhaFollow struct {
	BaseModel
	BuddhaId uint64 `json:"buddha_id"`
	UserId uint64 `json:"user_id"`
	WatchedAt int64 `json:"watched_at"`
	ExpireAt int64 `json:"expired_at"`
}

func (BuddhaFollow) TableName() string  {
	return "buddha_follows"
}