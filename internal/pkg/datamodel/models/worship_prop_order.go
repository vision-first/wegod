package models

type WorshipPropOrder struct {
	BaseModel
	UserId uint64
	PropId uint64
	PropAvailableDurationSnapshot int64
	Sn string
	Status int
}