package ws

import "github.com/zackradisic/youtube-rooms/internal/room"

// ClientMessage represents a message sent by a WebSocket client
type ClientMessage struct {
	Action string      `json:"action"`
	Data   interface{} `json:"data"`
	Client *Client
}

// HubMessage is a message sent by the hub to a WebSocket client
type HubMessage struct {
	JSON       []byte
	room       *room.Room
	recipients []*room.User
}

// NewHubMessage returns a hub message
func NewHubMessage(jsonData []byte, room *room.Room, recipients []*room.User) *HubMessage {

	if room != nil {
		return &HubMessage{
			JSON: jsonData,
			room: room,
		}
	} else if recipients != nil {
		return &HubMessage{
			JSON:       jsonData,
			recipients: recipients,
		}
	}

	return nil
}
