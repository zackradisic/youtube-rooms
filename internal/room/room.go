package room

import (
	"sync"

	"github.com/zackradisic/youtube-rooms/internal/models"
)

// Room manages a YouTube room
type Room struct {
	Model *models.Room
	Mux   sync.Mutex
	Users []*models.User
}

// NewRoom returns a new room
func NewRoom(model *models.Room) *Room {
	return &Room{
		Model: model,
		Users: make([]*models.User, 0),
	}
}
