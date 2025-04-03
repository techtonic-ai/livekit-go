package models

import (
	"gorm.io/gorm"
)

type Room struct {
	gorm.Model
	RoomID     string `json:"room_id" gorm:"unique;not null"`
	IsOccupied bool   `json:"is_occupied" gorm:"default:false"`
	UserID     string `json:"user_id" gorm:"not null"`
}
