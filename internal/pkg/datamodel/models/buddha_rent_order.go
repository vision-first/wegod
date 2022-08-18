package models

type BuddhaRentOrder struct {
	BaseModel
	UserId uint64 `json:"user_id"`
	BuddhaId uint64 `json:"buddha_id"`
	RentPackageId uint64 `json:"rent_package_id"`
	RentPackagePriceSnapshot uint32 `json:"rent_package_price_snapshot"`
	RentPackageAvailableDurationSnapshot uint32 `json:"rent_package_available_duration_snapshot"`
	RentPackageNameSnapshot string
	RentPackageDescSnapshot string
	PayedAt int64 `json:"payed_at"`
	EndRentAt int64 `json:"end_rent_at"`
	PaymentFlowId uint64 `json:"payment_flow_id"`
	Status uint `json:"status"`
	Sn string `json:"sn"`
	Money uint32
	Remark string
}
