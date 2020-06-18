package models

import "github.com/jinzhu/gorm"

// User is identified by their discord ID, and we keep track of their last known discord
// username and discriminator (i.e. bob#9832)
type User struct {
	gorm.Model
	DiscordID                string `gorm:"UNIQUE;UNIQUE_INDEX;NOT NULL"`
	LastDiscordUsername      string `gorm:"type:VARCHAR(32)"`
	LastDiscordDiscriminator string `gorm:"type:VARCHAR(4)"`
}
