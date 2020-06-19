package server

import (
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql" // MySQL driver
	"github.com/zackradisic/youtube-rooms/internal/models"
)

func (s *Server) setupDB() (*gorm.DB, error) {
	db, err := gorm.Open("mysql", "root@/youtube_rooms?charset=utf8mb4&parseTime=True&loc=Local")
	if err != nil {
		return nil, err
	}

	// Be careful! AutoMigrate() won't update existing columns or delete them!
	user, userAuth, room, video := &models.User{}, &models.UserAuth{}, &models.Room{}, &models.Video{}
	db.AutoMigrate(user, userAuth, room, video)
	return db, nil
}

func (s *Server) createUser(userInfo *discordUserInfoResponse, authToken *AuthToken) (*models.User, error) {
	user := &models.User{}
	if err := s.DB.Where(&models.User{DiscordID: userInfo.ID}).First(&user).Error; err != nil {
		if gorm.IsRecordNotFoundError(err) {
			user.DiscordID = userInfo.ID
			user.LastDiscordUsername = userInfo.Username
			user.LastDiscordDiscriminator = userInfo.Discriminator
			if err = s.DB.Create(user).Error; err != nil {
				return nil, err
			}
		} else {
			return nil, err
		}
	}

	return user, nil
}
