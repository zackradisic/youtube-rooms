// Copyright 2013 The Gorilla WebSocket Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package ws

import (
	"encoding/json"
	"fmt"

	"github.com/zackradisic/youtube-rooms/internal/room"
)

// Hub maintains the set of active clients and broadcasts messages to the
// clients.
type Hub struct {
	// Registered clients.
	clients map[*Client]bool

	users map[*room.User]*Client

	// Inbound messages from the clients.
	inbound chan *ClientMessage

	// Outbound messages to be sent to clients
	outbound chan *HubMessage

	// Register requests from the clients.
	register chan *Client

	// Unregister requests from clients.
	unregister chan *Client

	roomManager   *room.Manager
	actionInvoker *ActionInvoker
}

// NewHub creates a new Hub
func NewHub(roomManager *room.Manager) *Hub {
	outbound := make(chan *HubMessage)
	return &Hub{
		inbound:       make(chan *ClientMessage),
		outbound:      outbound,
		register:      make(chan *Client),
		unregister:    make(chan *Client),
		clients:       make(map[*Client]bool),
		users:         make(map[*room.User]*Client),
		roomManager:   roomManager,
		actionInvoker: NewActionInvoker(),
	}
}

// Run runs the Hub
func (h *Hub) Run() {
	for {
		select {
		case client := <-h.register:
			h.clients[client] = true
			h.users[client.user] = client
			fmt.Println("Registered client: " + client.user.Model.DiscordID)
		case client := <-h.unregister:
			if _, ok := h.clients[client]; ok {
				h.removeClient(client)
			}
		case message := <-h.inbound:
			fmt.Printf("Received message from client: (%s)\n", message.Client.user.DiscordHandle())
			err := h.actionInvoker.InvokeAction(message, h.outbound)
			if err != nil {
				fmt.Println(err)
			}
		case message := <-h.outbound:
			fmt.Println("Received outbound message")
			h.broadcastMessage(message)
		}
	}
}

func (h *Hub) broadcastMessage(message *HubMessage) {

	if message.recipients != nil {
		for _, recipient := range message.recipients {
			h.send(message.JSON, h.users[recipient])
		}
	} else if message.room != nil {
		for user, client := range h.users {
			if user.CurrentRoom.Model.Name == message.room.Model.Name {
				h.send(message.JSON, client)
			}
		}
	} else {
		fmt.Println("broadcastMessage() -> Received message with no recipients or room specified")
	}
}

func (h *Hub) send(msg []byte, client *Client) {
	if client == nil {
		return
	}

	select {
	case client.send <- msg:
	default:
		h.removeClient(client)
	}
}

func (h *Hub) removeClient(client *Client) {
	r := client.user.CurrentRoom
	client.user.CurrentRoom.RemoveUser(client.user)
	delete(h.clients, client)
	delete(h.users, client.user)
	close(client.send)

	// This is probably not the best way to manually broadcast a HubMessage but it is the quickest
	// solution with the current actions implementation.
	//
	// The actions system should have an additional layer of abstraction so that actions can be
	// invoked irrespective of clients: decouple actions from the client
	users := getUsersJSON(r)
	data, err := json.Marshal(&users)
	if err != nil {
		return
	}
	msg := NewHubMessage(data, r, nil)
	h.broadcastMessage(msg)
}
