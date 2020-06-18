package models

import "github.com/jinzhu/gorm"

// Room is a room a user can create, it can optionally be secured with a passowrd
type Room struct {
	gorm.Model
	Owner          uint   `gorm:"foreignkey:user_auth_user_fk_user_id;UNIQUE_INDEX;NOT_NULL"`
	HashedPassword string `gorm:"type:CHAR(76)"`
}
