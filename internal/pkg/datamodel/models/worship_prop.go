package models

type WorshipProp struct {
	BaseModel
	Name string `json:"name" gorm:"type:varchar"`
	Image string `json:"image" gorm:"type:varchar"`
	Price uint32 `json:"price"`
	ShelfStatus       uint `json:"shelf_status"`
	AvailableDuration int64 `json:"available_duration"`
}

func (w *WorshipProp) IsFree() bool {
	return w.Price == 0
}