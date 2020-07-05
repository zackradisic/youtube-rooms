package room

import (
	"fmt"
	"sync"

	"github.com/zackradisic/youtube-rooms/internal/models"
)

// Room manages a YouTube room
type Room struct {
	Model     *models.Room
	mux       sync.Mutex
	Users     []*User
	Current   *Video
	IsPlaying bool
}

// NewRoom returns a new room
func NewRoom(model *models.Room) *Room {
	return &Room{
		Model: model,
		Users: make([]*User, 0),
	}
}

// GetUsers returns the Users currently in this room
func (r *Room) GetUsers() []*User {
	r.mux.Lock()
	defer r.mux.Unlock()
	return r.Users
}

// SetIsPlaying sets whether or not the room is playing the video
func (r *Room) SetIsPlaying(isPlaying bool) {
	r.mux.Lock()
	defer r.mux.Unlock()
	r.IsPlaying = isPlaying
}

// SetCurrentVideo sets this room's current video
func (r *Room) SetCurrentVideo(video *Video) {
	r.mux.Lock()
	defer r.mux.Unlock()
	r.Current = video
}

// AddUser adds a user to the room
func (r *Room) AddUser(user *User) error {
	r.mux.Lock()
	defer r.mux.Unlock()
	i := r.getUserIndex(user)
	if i != -1 {
		return fmt.Errorf("already in the room")
	}
	r.Users = append(r.Users, user)
	user.CurrentRoom = r
	return nil
}

// RemoveUser removes a user from the rooms
func (r *Room) RemoveUser(user *User) {
	r.mux.Lock()
	defer r.mux.Unlock()
	i := r.getUserIndex(user)
	if i == -1 {
		return
	}
	r.Users = append(r.Users[:i], r.Users[i+1:]...)
	user.CurrentRoom = nil
}

// HasUser returns true if the user is in this room, false if not
func (r *Room) HasUser(user *User) bool {
	r.mux.Lock()
	defer r.mux.Unlock()
	for _, u := range r.Users {
		if u.Model.DiscordID == user.Model.DiscordID {
			return true
		}
	}

	return false
}

func (r *Room) getUserIndex(user *User) int {
	for i, u := range r.Users {
		if u.Model.DiscordID == user.Model.DiscordID {
			return i
		}
	}

	return -1
}
