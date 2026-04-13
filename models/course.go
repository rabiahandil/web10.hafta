package models

import (
	"gorm.io/gorm"
)

type Course struct {
	gorm.Model
	Title       string   `json:"title"`
	Description string   `json:"description"`
	Category    string   `json:"category"`
	TeacherID   uint     `json:"teacher_id"`
	Teacher     User     `gorm:"foreignKey:TeacherID" json:"teacher"`
	Lessons     []Lesson `json:"lessons"`
}
