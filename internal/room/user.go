package room

import "github.com/zackradisic/youtube-rooms/internal/models"

// User represents a user in the room
type User struct {
	Model *models.User
	*models.UserAuth
	CurrentRoom *Room
}

// NewUser creates a new user
func NewUser(user *models.User, userAuth *models.UserAuth) *User {
	newUser := &User{}
	newUser.Model = user
	newUser.UserAuth = userAuth
	return newUser
}
