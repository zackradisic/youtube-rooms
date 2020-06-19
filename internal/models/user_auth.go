package models

import (
	"time"

	"github.com/jinzhu/gorm"
)

// UserAuth contains the Discord OAuth credentials for the user identified by foreign key,
type UserAuth struct {
	gorm.Model
	AccessToken  string    `gorm:"type:VARCHAR(30);index"`
	RefreshToken string    `gorm:"type:VARCHAR(30);index"`
	ExpiresAt    time.Time `gorm:"index"`
	UserID       uint      `gorm:"UNIQUE_INDEX;NOT_NULL"`
	User         User
}

// TableName sets the UserAuth model's table name as user_auth
func (UserAuth) TableName() string {
	return "user_auth"
}
