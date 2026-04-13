package models

import (
	"gorm.io/gorm"
)

type Quiz struct {
	gorm.Model
	LessonID  uint       `json:"lesson_id"`
	Questions []Question `json:"questions"`
}

type Question struct {
	gorm.Model
	QuizID  uint   `json:"quiz_id"`
	Text    string `json:"text"`
	OptionA string `json:"option_a"`
	OptionB string `json:"option_b"`
	OptionC string `json:"option_c"`
	OptionD string `json:"option_d"`
	Correct string `json:"correct"` // A, B, C, D
}

type QuizResult struct {
	gorm.Model
	UserID  uint    `json:"user_id"`
	QuizID  uint    `json:"quiz_id"`
	Score   int     `json:"score"`
	Total   int     `json:"total"`
	Percent float64 `json:"percent"`
}
