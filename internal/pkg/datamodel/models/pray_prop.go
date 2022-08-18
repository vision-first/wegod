package models

import (
	"github.com/995933447/dbdriverutil/field"
)

type PrayProp struct {
	BaseModel
	Name string `json:"name" gorm:"type:varchar"`
	Price uint32 `json:"price"`
	Remark string `json:"remark" gorm:"type:varchar"`
	ShelfStatus int `json:"shelf_status"`
	Extra field.Json `json:"extra"`
}
