package room

import "time"

// User represents a user in the room
type User struct {
	DiscordID            string
	DiscordUsername      string
	DiscordDiscriminator string
	AccessToken          string
	AccessTokenExpiry    *time.Time
}
