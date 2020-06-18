package server

import (
	"github.com/jinzhu/gorm"
	"github.com/zackradisic/youtube-rooms/internal/models"
)

func (s *Server) setupDB() (*gorm.DB, error) {
	// Using SQLite for development, may want to consider
	// changing to MySQL or Postgres later for production
	db, err := gorm.Open("sqlite3", "test.db")
	if err != nil {
		return nil, err
	}

	defer db.Close()

	// Be careful! AutoMigrate() won't update existing columns or delete them!
	db.AutoMigrate(&models.User{}, &models.UserAuth{})

	return db, nil
}
