package models

type Post struct {
	BaseModel
	Title string `json:"title" gorm:"type:varchar"`
	Content string `json:"content"`
}
