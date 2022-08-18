package models

type UserWorshipProp struct {
	BaseModel
	UserId uint64
	PropId uint64
	AvailableDuration int64
}