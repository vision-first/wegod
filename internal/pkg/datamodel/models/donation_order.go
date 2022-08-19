package models

type DonationScene int

const (
	DonationSceneIncense DonationScene = iota
	DonationSceneTriad
	DonationSceneFreeCaptive
	DonationSceneSealBook
)

type DonationOrder struct {
	BaseModel
	UserId uint64 `json:"user_id" gorm:"index"`
	BuddhaId uint64 `json:"buddha_id" gorm:"index"`
	Money uint32
	PayedAt int64
	DonationScene uint32
	Remark string `json:"remark" gorm:"type:varchar"`
	Sn string
	Status int
}
