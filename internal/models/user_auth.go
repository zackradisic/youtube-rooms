package models

import "github.com/jinzhu/gorm"

// UserAuth contains the Discord OAuth credentials for the user identified by foreign key
type UserAuth struct {
	gorm.Model
	AccessToken string `gorm:"type:VARHCAR(30);inde"`
}
