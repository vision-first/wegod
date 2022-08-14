package datamodels

type Buddha struct {
	BaseModel
	Name string `json:"name" gorm:"type:varchar"`
	Image string `json:"image" gorm:"type:varchar"`
	Sort uint32 `json:"sort"`
	Remark string `json:"remark" gorm:"type:varchar"`
}

func (Buddha) TableName() string  {
	return "buddhas"
}