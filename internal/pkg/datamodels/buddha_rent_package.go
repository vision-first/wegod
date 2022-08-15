package datamodels

type BuddhaRent struct {
	BaseModel
	BuddhaId uint64 `json:"buddha_id"`
	Name string `json:"name"`
	Desc string `json:"desc"`
	Price uint64 `json:"price"`
	AvailableDuration uint32 `json:"effective_time"`
}
