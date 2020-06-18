package models

import "github.com/jinzhu/gorm"

// UserAuth contains the Discord OAuth credentials for the user identified by foreign key,
type UserAuth struct {
	gorm.Model
	AccessToken  string `gorm:"type:VARCHAR(30);index"`
	RefreshToken string `gorm:"type:VARCHAR(30);index"`
	ExpiresIn    uint   `gorm:"type:UNSIGNED INT(10);index"`
	User         uint   `gorm:"foreignkey:user_auth_user_fk_user_id;UNIQUE_INDEX;NOT_NULL"`
}

// TableName sets the UserAuth model's table name as user_auth
func (UserAuth) TableName() string {
	return "user_auth"
}
