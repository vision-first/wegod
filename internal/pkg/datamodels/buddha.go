package datamodels

type Buddha struct {
	BaseModel
	Name string
	Image string
	Sort uint32
	Remark string
}

type UserRefBuddha struct {
	BaseModel
	BuddhaId uint64
	Uid uint64
}