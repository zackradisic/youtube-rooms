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
	JSON     []byte
	roomName string
	room     *room.Room
}

// NewHubMessage returns a hub message
func NewHubMessage(jsonData []byte, room *room.Room) *HubMessage {

	return &HubMessage{
		JSON:     jsonData,
		roomName: room.Model.Name,
		room:     room,
	}
}
