package models

type UserDonationDailyStat struct {
	BaseModel
	UserId uint64 `json:"user_id"`
	TotalPayedMoney uint32 `json:"total_payed_money"`
	TotalMoney uint32 `json:"total_money"`
	TotalOrderNum uint32 `json:"total_order_num"`
	TotalPayedOrderNum uint32 `json:"total_payed_order_num"`
	DateYmd int `json:"date_ymd"`
	CalculatedAt int64 `json:"calculated_at"`
}
