package models

type UserBuddhaRent struct {
	BaseModel
	UserId uint64
	BuddhaId uint64
	ExpireAt int64
}
