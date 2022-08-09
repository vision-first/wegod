package datamodels

const (
	DonationMotiveIncense= iota
	DonationMotiveTriad
	DonationMotiveFreeCaptive
	DonationMotiveSealBook
)

type UserDonationRecord struct {
	BaseModel
	UserId uint64 `json:"user_id" gorm:"index"`
	Money uint32
	PayedAt int64
	DonationMotive uint32
	Remark string `json:"remark" gorm:"type:varchar"`
}
