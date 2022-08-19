package models

const (
	OrderTypeBuddhaRent = 1
	OrderTypePrayProp = 2
	OrderTypeWorshipProp = 3
)

type PaymentFlow struct {
	BaseModel
	UserId uint64
	OrderId uint64
	OrderType uint32
}
