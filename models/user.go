package models

import (
	"gorm.io/gorm"
)

type User struct {
	gorm.Model
	Name   string `json:"name" gorm:"not null"`
	UserID string `json:"user_id" gorm:"unique;not null"`
}
