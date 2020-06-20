package room

import (
	"sync"

	"github.com/zackradisic/youtube-rooms/internal/models"
)

// Room manages a YouTube room
type Room struct {
	Model     *models.Room
	Mux       sync.Mutex
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

// SetIsPlaying sets whether or not the room is playing the video
func (r *Room) SetIsPlaying(isPlaying bool) {
	r.Mux.Lock()
	defer r.Mux.Unlock()
	r.IsPlaying = isPlaying
}

// SetCurrentVideo sets this room's current video
func (r *Room) SetCurrentVideo(video *Video) {
	r.Mux.Lock()
	defer r.Mux.Unlock()
	r.Current = video
}

// AddUser adds a user to the room
func (r *Room) AddUser(user *User) {
	r.Mux.Lock()
	defer r.Mux.Unlock()
	i := r.getUserIndex(user)
	if i != -1 {
		return
	}
	r.Users = append(r.Users, user)
	user.CurrentRoom = r
}

// RemoveUser removes a user from the rooms
func (r *Room) RemoveUser(user *User) {
	r.Mux.Lock()
	defer r.Mux.Unlock()
	i := r.getUserIndex(user)
	if i == -1 {
		return
	}
	r.Users = append(r.Users[:i], r.Users[i+1:]...)
	user.CurrentRoom = nil
}

// HasUser returns true if the user is in this room, false if not
func (r *Room) HasUser(user *User) bool {
	r.Mux.Lock()
	defer r.Mux.Unlock()
	for _, u := range r.Users {
		if u.Model.DiscordID == user.Model.DiscordID {
			return true
		}
	}

	return false
}

func (r *Room) getUserIndex(user *User) int {
	for i, u := range r.Users {
		if u.User.DiscordID == user.User.DiscordID {
			return i
		}
	}

	return -1
}
