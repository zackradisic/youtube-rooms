package room

import "github.com/zackradisic/youtube-rooms/internal/models"

// User represents a user in the room
type User struct {
	models.User
	models.UserAuth
	CurrentRoom *Room
}
