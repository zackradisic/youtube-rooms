package models

import "github.com/jinzhu/gorm"

// Room is a room a user can create, it can optionally be secured with a passowrd
type Room struct {
	gorm.Model
	Owner          User   `gorm:"foreignkey:OwnerID"`
	OwnerID        uint   `gorm:"UNIQUE_INDEX;NOT_NULL"`
	HashedPassword string `gorm:"type:CHAR(76)"`
	Name           string `gorm:"type:VARCHAR(36);UNIQUE_INDEX;NOT_NULL"`
}
