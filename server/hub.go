package server

import (
	"encoding/json"
	"log"
	"net/http"
	"sync"

	"github.com/gorilla/websocket"
)

type Hub struct {
	clients    map[string]*Client
	register   chan *Client
	unregister chan *Client
	handlers   map[string]HandleFunc
	mu         sync.RWMutex
}

func NewHub() *Hub {
	return &Hub{
		clients:    make(map[string]*Client),
		register:   make(chan *Client),
		unregister: make(chan *Client),
		handlers:   make(map[string]HandleFunc),
	}
}

func (h *Hub) Run() {
	for {
		select {
		case client := <-h.register:
			h.mu.Lock()
			h.clients[client.ID] = client
			h.mu.Unlock()

		case client := <-h.unregister:
			h.mu.Lock()
			if _, ok := h.clients[client.ID]; ok {
				delete(h.clients, client.ID)
				client.Close()
			}
			h.mu.Unlock()
		}
	}
}

func (h *Hub) On(event string, handler HandleFunc) {
	h.mu.Lock()
	defer h.mu.Unlock()
	h.handlers[event] = handler
}

func (h *Hub) dispatch(client *Client, payload []byte) {
	var packet Packet
	if err := json.Unmarshal(payload, &packet); err != nil {
		_ = client.Emit("error", map[string]string{"message": "invalid packet format"})
		return
	}
	h.mu.RLock()
	handler, exists := h.handlers[packet.Event]
	h.mu.RUnlock()

	if !exists {
		log.Printf("unknown event: %s", packet.Event)
		return
	}
	handler(client, packet.Data)
}

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

func Upgrade(w http.ResponseWriter, r *http.Request) (*websocket.Conn, error) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		return nil, err
	}
	return conn, nil
}
