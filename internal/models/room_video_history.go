package models

import "github.com/jinzhu/gorm"

// RoomVideo represents a video in a room's video history
type RoomVideo struct {
	gorm.Model
	Room        Room
	RoomID      uint  `gorm:"INDEX;NOT NULL"`
	Requester   User  `gorm:"foreignkey:RequesterID"`
	RequesterID uint  `gorm:"INDEX"`
	Video       Video `gorm:"foreignkey:VideoID"`
	VideoID     uint  `gorm:"INDEX"`
}
