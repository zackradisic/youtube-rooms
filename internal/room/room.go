package room

import (
	"sync"

	"github.com/zackradisic/youtube-rooms/internal/models"
)

// Room manages a YouTube room
type Room struct {
	Model   *models.Room
	Mux     sync.Mutex
	Users   []*User
	Current *Video
}

// NewRoom returns a new room
func NewRoom(model *models.Room) *Room {
	return &Room{
		Model: model,
		Users: make([]*User, 0),
	}
}

// AddUser adds a user to the room
func (r *Room) AddUser(user *User) {
	i := r.getUserIndex(user)
	if i != -1 {
		return
	}
	r.Users = append(r.Users, user)
}

// RemoveUser removes a user from the rooms
func (r *Room) RemoveUser(user *User) {
	i := r.getUserIndex(user)
	if i == -1 {
		return
	}
	r.Users = append(r.Users[:i], r.Users[i+1:]...)
}

func (r *Room) getUserIndex(user *User) int {
	for i, u := range r.Users {
		if u.User.DiscordID == user.User.DiscordID {
			return i
		}
	}

	return -1
}
