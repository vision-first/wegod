package datamodels

type BuddhaWorshipProp struct {
	BaseModel
	Name string `json:"name" gorm:"type:varchar"`
	Image string `json:"image" gorm:"type:varchar"`
	Price uint32 `json:"price"`
	ShelfStatus uint32 `json:"shelf_status"`
}