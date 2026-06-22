package server

import (
	"encoding/json"
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

func TestHub_OnAndDispatch_Success(t *testing.T) {
	hub := NewHub()
	client := NewClient(nil, hub)

	var isCalled bool
	var receivedMessage string

	// Đăng ký sự kiện "chat:send"
	hub.On("chat:send", func(c *Client, data json.RawMessage) {
		isCalled = true
		var msg string
		if err := json.Unmarshal(data, &msg); err != nil {
			t.Errorf("Lỗi giải mã JSON data: %v", err)
		}
		receivedMessage = msg
	})

	// Giả lập client gửi thông điệp JSON tương ứng lên server
	rawPayload := []byte(`{"event":"chat:send","data":"xin chao gocket"}`)
	hub.dispatch(client, rawPayload)

	if !isCalled {
		t.Fatal("Lỗi: Event handler đã không được gọi khi dispatch sự kiện")
	}

	if receivedMessage != "xin chao gocket" {
		t.Errorf("Lỗi: Nhận sai nội dung sự kiện. Mong muốn 'xin chao gocket', nhận được '%s'", receivedMessage)
	}
}

func TestHub_Dispatch_InvalidJSON(t *testing.T) {
	hub := NewHub()
	client := NewClient(nil, hub)

	// Chạy goroutine kiểm tra xem client có nhận được gói tin sự kiện "error" hay không
	done := make(chan struct{})
	go func() {
		defer close(done)
		msg, ok := <-client.Send
		if !ok {
			t.Error("Kênh Send đã bị đóng")
			return
		}

		var packet Packet
		if err := json.Unmarshal(msg, &packet); err != nil {
			t.Errorf("Không giải mã được phản hồi lỗi: %v", err)
			return
		}

		if packet.Event != "error" {
			t.Errorf("Mong muốn nhận sự kiện loại 'error', nhận được '%s'", packet.Event)
		}
	}()

	// Gửi JSON không hợp lệ
	rawPayload := []byte(`{invalid-json}`)
	hub.dispatch(client, rawPayload)

	<-done
}