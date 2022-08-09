package datamodels

type PrayProp struct {
	BaseModel
	Name string `json:"name" gorm:"type:varchar"`
	Price uint32 `json:"price"`
	Remark string `json:"remark" gorm:"type:varchar"`
}
