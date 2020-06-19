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
