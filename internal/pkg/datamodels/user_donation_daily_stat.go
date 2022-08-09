package datamodels

type UserDonationDailyStat struct {
	BaseModel
	UserId uint64
	TotalPayedMoney uint32
	dateYmd uint32
	CalculatedAt uint32
}
