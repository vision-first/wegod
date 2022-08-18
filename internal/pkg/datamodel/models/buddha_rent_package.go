package models

type BuddhaRentPackage struct {
	BaseModel
	BuddhaId uint64 `json:"buddha_id"`
	Name string `json:"name"`
	Desc string `json:"desc"`
	Price uint32 `json:"price"`
	AvailableDuration uint32 `json:"effective_time"`
	ShelfStatus int `json:"shelf_status"`
	Sort uint32 `json:"sort"`
}
