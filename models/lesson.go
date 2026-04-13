package models

import (
	"gorm.io/gorm"
)

type Lesson struct {
	gorm.Model
	Title    string `json:"title"`
	Content  string `json:"content"`
	VideoURL string `json:"video_url"`
	Order    int    `json:"order"`
	CourseID uint   `json:"course_id"`
}
