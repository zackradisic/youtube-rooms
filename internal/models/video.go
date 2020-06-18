package models

import "github.com/jinzhu/gorm"

// Video represents a YouTube video that is played in a room
type Video struct {
	gorm.Model
	Room      uint   `gorm:"foreignkey:videos_room_fk_rooms_id;INDEX;NOT_NULL"`
	Title     string `gorm:"type:VARCHAR(100);INDEX;NOT_NULL"`
	YouTubeID string `gorm:"type:VARCHAR(11);UNIQUE_INDEX;NOT_NULL"`
	Requester uint   `gorm:"foreignkey:videos_requester_fk_users_id;INDEX"`
}
