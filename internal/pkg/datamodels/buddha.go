package datamodels

type Buddha struct {
	BaseModel
	Name string
	Image string
	Sort uint32
	Remark string
}

type UserWatchedBuddha struct {
	BaseModel
	BuddhaId uint64
	Uid uint64
	LastWatchedAt int64
}