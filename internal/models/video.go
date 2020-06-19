package models

import "github.com/jinzhu/gorm"

// Video represents a YouTube video that is played in a room
type Video struct {
	gorm.Model
	Room        Room
	RoomID      uint   `gorm:"INDEX;NOT NULL"`
	Title       string `gorm:"type:VARCHAR(100);INDEX;NOT NULL"`
	YoutubeID   string `gorm:"type:VARCHAR(11);UNIQUE_INDEX;NOT NULL"`
	Requester   User   `gorm:"foreignkey:RequesterID"`
	RequesterID uint   `gorm:"INDEX"`
}
