package models

const (
	GenderNil = iota
	GenderMan
	GenderWoman
)

type User struct {
	BaseModel
	Phone string `json:"phone" gorm:"type:varchar"`
	HashedPassword string `json:"password" gorm:"type:varchar"`
	NickName string `json:"nick_name" gorm:"type:varchar"`
	Avatar string `json:"avatar" gorm:"type:varchar"`
	Desc string `json:"desc"`
	Gender uint8 `json:"gender"`
	ChannelId uint64 `json:"channel_id"`
}

