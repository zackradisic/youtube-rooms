package models

import "github.com/jinzhu/gorm"

// UserAuth contains the Discord OAuth credentials for the user identified by foreign key,
type UserAuth struct {
	gorm.Model
	AccessToken  string `gorm:"type:VARCHAR(30);index"`
	RefreshToken string `gorm:"type:VARCHAR(30);index"`
	ExpiresIn    uint   `gorm:"type:UNSIGNED INT(10);index"`
}
