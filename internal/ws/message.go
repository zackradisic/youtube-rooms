package ws

// Message represents a message sent by a WebSocket client
type Message struct {
	Action   string      `json:"action"`
	RoomName string      `json:"roomName"`
	Data     interface{} `json:"data"`
}
