package models

import (
	"gorm.io/gorm"
)

type User struct {
	gorm.Model
	Name     string `json:"name"`
	Email    string `gorm:"uniqueIndex;not null" json:"email"`
	Password string `json:"-"`
	Role     string `gorm:"default:'student'" json:"role"` // student, teacher
}
