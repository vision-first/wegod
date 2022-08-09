package datamodels

type BuddhaFollow struct {
	BaseModel
	BuddhaId uint64 `json:"buddha_id"`
	UserId uint64 `json:"user_id"`
	WatchedAt int64 `json:"watched_at"`
}
