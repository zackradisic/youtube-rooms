package room

import (
	"fmt"

	"github.com/zackradisic/youtube-rooms/internal/models"
)

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

// DiscordHandle returns the user's username and discriminator as it appears on discord
// Ex: jeremy#5932
func (u *User) DiscordHandle() string {
	return fmt.Sprintf("%s#%s", u.Model.LastDiscordUsername, u.Model.LastDiscordDiscriminator)
}
