package server

import (
	"testing"
	"time"
)

func TestHub_RegisterAndUnregister(t *testing.T) {
	hub := NewHub()
	go hub.Run()

	client := &Client{
		ID:   "test-client-unique-id",
		Send: make(chan []byte, 10),
		hub:  hub,
	}

	hub.register <- client
	time.Sleep(50 * time.Millisecond)

	hub.mu.RLock()
	_, exists := hub.clients[client.ID]
	hub.mu.RUnlock()

	if !exists {
		t.Fatal("Lỗi: Client đáng lẽ phải được đăng ký thành công vào Hub")
	}

	hub.unregister <- client
	time.Sleep(50 * time.Millisecond)

	hub.mu.RLock()
	_, exists = hub.clients[client.ID]
	hub.mu.RUnlock()

	if exists {
		t.Fatal("Lỗi: Client vẫn tồn tại trong Hub sau khi thực hiện unregister")
	}
}