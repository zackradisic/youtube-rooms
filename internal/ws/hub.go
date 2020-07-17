// Copyright 2013 The Gorilla WebSocket Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package ws

import (
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
			fmt.Println("Registed client: " + client.user.Model.DiscordID)
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

	for user := range h.users {
		if user.CurrentRoom.Model.Name == message.roomName {
			select {
			case h.users[user].send <- message.JSON:
				fmt.Println(string(message.JSON))
			default:
				h.removeClient(h.users[user])
			}

		}
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
	userProxy := &room.User{
		CurrentRoom: r,
	}
	clientProxy := &Client{
		user: userProxy,
	}
	msg, _ := getUsersAction(nil, clientProxy)
	h.broadcastMessage(msg)
}
