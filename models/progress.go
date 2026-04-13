package models

import (
	"gorm.io/gorm"
)

type Progress struct {
	gorm.Model
	UserID   uint `gorm:"uniqueIndex:idx_user_lesson" json:"user_id"`
	LessonID uint `gorm:"uniqueIndex:idx_user_lesson" json:"lesson_id"`
	CourseID uint `json:"course_id"`
}
