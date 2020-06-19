package room

import "github.com/jinzhu/gorm"

// Manager manages rooms
type Manager struct {
	rooms map[string]*Room
	DB    *gorm.DB
}

// NewManager returns a new manager
func NewManager(db *gorm.DB) *Manager {
	return &Manager{
		rooms: make(map[string]*Room),
		DB:    db,
	}
}
