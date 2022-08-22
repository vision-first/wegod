package dtos

type Buddha struct {
	Id uint64 `json:"id"`
	Name string `json:"name"`
	Image string `json:"image"`
	Sort uint32 `json:"sort"`
}

type PageBuddhaResp struct {
	List []*Buddha `json:"list"`
	PageQueryResp
}

type WatchBuddhaReq struct {
	BuddhaId uint64 `json:"buddha_id"`
}

type WatchBuddhaResp struct {
}

type UnwatchBuddhaReq struct {
	BuddhaId uint64 `json:"buddha_id"`
}

type UnwatchBuddhaResp struct {
}

type PageUserWatchedBuddhasResp struct {
	List []*Buddha
	PageQueryResp
}

type CreateBuddhaRentOrderReq struct {
	RentPackageId uint64
}

type BuddhaRentOrder struct {
	OrderSn string
	BuddhaId uint64
	RentPackageId uint64
	Price uint32
	RentPackageName string
	RentPackageDesc string
	Status uint
}

type CreateBuddhaRentOrderResp struct {
	Order *BuddhaRentOrder `json:"order"`
}

type BuddhaRentPackage struct {
	Id uint64 `json:"id"`
	Name string `json:"name"`
	Desc string `json:"desc"`
	Price uint32 `json:"price"`
	AvailableDuration uint32 `json:"effective_time"`
	ShelfStatus int `json:"shelf_status"`
	BuddhaId uint64 `json:"buddha_id"`
}

type PageBuddhaRentPackagesResp struct {
	List []*BuddhaRentPackage `json:"list"`
	PageQueryResp
}

type PrayToBuddhaReq struct {
	BuddhaId uint64 `json:"buddha_id"`
	PrayPropIds []uint64 `json:"pray_prop_ids"`
	Content string	`json:"content"`
}

type PrayToBuddhaResp struct {
}

type WorshipToBuddhaReq struct {
	BuddhaId uint64 `json:"buddha_id"`
	WorshipPropId uint64 `json:"worship_prop_id"`
	ConsumeUserWorshipPropId uint64 `json:"consume_user_worship_prop_id"`
}

type WorshipToBuddhaResp struct {
}