package models

const (
	PaymentChannelWechat = iota
)

type UserPaymentFlow struct {
	BaseModel
	UserId uint64
	Money uint32
	PayedAt int64
	Channel int
}
