package room

import (
	"fmt"
	"sync"

	"github.com/jinzhu/gorm"
	"github.com/zackradisic/youtube-rooms/internal/models"
)

// Manager manages rooms and users
type Manager struct {
	rooms map[string]*Room
	DB    *gorm.DB
	mux   sync.Mutex
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
	m.mux.Lock()
	defer m.mux.Unlock()
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

	// Cache it
	newRoom := NewRoom(room)
	m.rooms[room.Name] = newRoom

	return newRoom, nil
}

// RemoveUser removes a user from the room it is in, if the user is not in a room it does
// nothing
func (m *Manager) RemoveUser(user *User) {
	m.mux.Lock()
	defer m.mux.Unlock()
	for _, room := range m.rooms {
		if room.Model.Name == user.CurrentRoom.Model.Name {
			room.RemoveUser(user)
		}
	}
}
