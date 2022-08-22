package models

type Post struct {
	BaseModel
	Title string `json:"title" gorm:"type:varchar"`
	Content string `json:"content"`
	Sort uint32 `json:"sort"`
	CategoryId uint64 `json:"category_id"`
}
