package server

import (
	"time"

	"github.com/google/uuid"
	"github.com/gorilla/websocket"
)

type Client struct {
	ID 			string
	Conn 		*websocket.Conn
	Send 		chan []byte
	hub 		*Hub
	rooms		map[string]bool
	connectedAt time.Time
}

func NewClient(conn *websocket.Conn, hub *Hub) *Client {
	return &Client{
		ID: 			uuid.New().String(),
		Conn: 			conn,
		Send: 			make(chan []byte, 256),
		hub: 			hub, 
		rooms: 			make(map[string]bool),
		connectedAt: 	time.Now(),
	}
}