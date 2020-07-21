package server

import (
	"os"
	"time"

	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql" // MySQL driver
	"github.com/zackradisic/youtube-rooms/internal/models"
)

func (s *Server) setupDB() (*gorm.DB, error) {
	db, err := gorm.Open("mysql", os.Getenv("DB_CONNECTION_STRING"))
	if err != nil {
		return nil, err
	}

	// Be careful! AutoMigrate() won't update existing columns or delete them!
	user, userAuth, room, video := &models.User{}, &models.UserAuth{}, &models.Room{}, &models.Video{}
	db.AutoMigrate(user, userAuth, room, video)
	db.Model(&models.UserAuth{}).AddForeignKey("user_id", "users(id)", "CASCADE", "CASCADE")
	db.Model(&models.Room{}).AddForeignKey("owner_id", "users(id)", "CASCADE", "CASCADE")
	db.Model(&models.Video{}).AddForeignKey("room_id", "rooms(id)", "CASCADE", "CASCADE")
	db.Model(&models.Video{}).AddForeignKey("requester_id", "users(id)", "CASCADE", "CASCADE")

	return db, nil
}

func (s *Server) getRooms(name string) (*[]models.Room, error) {
	rooms := []models.Room{}
	if name == "" {
		if err := s.DB.Find(&rooms).Error; err != nil {
			return nil, err
		}
	} else {
		if err := s.DB.Find(&rooms, "name = ?", name).Error; err != nil {
			return nil, err
		}
	}

	return &rooms, nil
}

func (s *Server) createUser(userInfo *discordUserInfoResponse, authToken *AuthToken) (*models.User, error) {

	user := &models.User{}
	auth := &models.UserAuth{}
	auth.AccessToken = authToken.AccessToken
	auth.RefreshToken = authToken.RefreshToken
	auth.ExpiresAt = time.Now().Add(time.Second * time.Duration(int64(authToken.ExpiresIn)))

	tx := s.DB.Begin()
	if err := tx.Where(&models.User{DiscordID: userInfo.ID}).First(&user).Error; err != nil {
		if gorm.IsRecordNotFoundError(err) {
			user.DiscordID = userInfo.ID
			user.LastDiscordUsername = userInfo.Username
			user.LastDiscordDiscriminator = userInfo.Discriminator

			if err = tx.Create(user).Error; err != nil {
				tx.Rollback()
				return nil, err
			}

			auth.UserID = user.ID
			if err := tx.Create(auth).Error; err != nil {
				tx.Rollback()
				return nil, err
			}

			tx.Commit()
		} else {
			tx.Rollback()
			return nil, err
		}

		return user, nil
	}

	auth.UserID = user.ID
	if err := tx.Model(&models.UserAuth{}).Where("user_id = ?", auth.UserID).Update(auth).Error; err != nil {
		tx.Rollback()
		return nil, err
	}

	if err := tx.Model(user).Update("last_discord_username", userInfo.Username).Error; err != nil {
		tx.Rollback()
		return nil, err
	}

	if err := tx.Model(user).Update("last_discord_discriminator", userInfo.Discriminator).Error; err != nil {
		tx.Rollback()
		return nil, err
	}

	tx.Commit()
	return user, nil
}
