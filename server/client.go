package server

import (
	"sync"
	"time"

	"github.com/google/uuid"
	"github.com/gorilla/websocket"
)

type Client struct {
	ID          string
	Conn        *websocket.Conn
	Send        chan []byte
	hub         *Hub
	rooms       map[string]bool
	connectedAt time.Time
	once        sync.Once
}

func NewClient(conn *websocket.Conn, hub *Hub) *Client {
	return &Client{
		ID:          uuid.New().String(),
		Conn:        conn,
		Send:        make(chan []byte, 256),
		hub:         hub,
		rooms:       make(map[string]bool),
		connectedAt: time.Now(),
	}
}

func (c *Client) Close() {
	c.once.Do(func() {
		close(c.Send)
		if c.Conn != nil {
			c.Conn.Close()
		}
	})
}

func (c *Client) readPump() {
	defer func() {
		c.hub.unregister <- c
		c.Conn.Close()
	}()

	for {
		_, message, err := c.Conn.ReadMessage()
		if err != nil {
			break
		}
		c.hub.dispatch(c, message)
	}
}

func (c *Client) writePump() {
	defer func() {
		if c.Conn != nil {
			c.Conn.Close()
		}
	}()

	for {
		select {
		case message, ok := <-c.Send:
			if !ok {
				if c.Conn != nil {
					_ = c.Conn.WriteMessage(websocket.CloseMessage, []byte{})
				}
				return
			}
			err := c.Conn.WriteMessage(websocket.TextMessage, message)
			if err != nil {
				return
			}
		}
	}
}