package models

type Buddha struct {
	BaseModel
	Name string `json:"name" gorm:"type:varchar"`
	Image string `json:"image" gorm:"type:varchar"`
	Sort uint32 `json:"sort"`
	Remark string `json:"remark" gorm:"type:varchar"`
	IsFreeRent bool `json:"is_free_rent"`
}

func (Buddha) TableName() string  {
	return "buddhas"
}