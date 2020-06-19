package room

import (
	"fmt"

	"github.com/jinzhu/gorm"
	"github.com/zackradisic/youtube-rooms/internal/models"
)

// Manager manages rooms and users
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

// GetRoom returns a room specified by the given name and caches it
func (m *Manager) GetRoom(name string) (*Room, error) {

	// First check the internal cache to see if the room exists
	if room, ok := m.rooms[name]; ok {
		return room, nil
	}

	// If not found lets check the DB
	room := &models.Room{
		Name: name,
	}
	if err := m.DB.First(room).Error; err != nil {
		return nil, fmt.Errorf("could not find that room (%s)", name)
	}

	newRoom := NewRoom(room)
	m.rooms[room.Name] = newRoom

	return newRoom, nil
}

// AddUser adds a user to a room, if the user is already in the room it does
// nothing
func (m *Manager) AddUser(room *Room, user *User) {
	room.AddUser(user)
}

// RemoveUser removes a user from a room, if the user is not in the room it does
// nothing
func (m *Manager) RemoveUser(room *Room, user *User) {
	room.RemoveUser(user)
}
