package models

type ShopProductCategory struct {
	BaseModel `json:"base_model"`
	Name string `json:"name" gorm:"type:varchar"`
	Level uint8 `json:"level"`
	Sort uint16 `json:"sort"`
}
