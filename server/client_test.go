package server

import (
	"encoding/json"
	"testing"
)

func TestClient_Emit_SuccessAndFailure(t *testing.T) {
	client := NewClient(nil, nil)

	// Dữ liệu giả lập để gửi thử
	testPayload := map[string]string{"message": "hello gocket"}

	// 1. Chạy bất đồng bộ để đọc từ kênh Send của Client và kiểm tra Packet
	go func() {
		msg, ok := <-client.Send
		if !ok {
			return
		}

		var packet Packet
		if err := json.Unmarshal(msg, &packet); err != nil {
			t.Errorf("Dữ liệu nhận từ kênh Send không giải mã được thành Packet: %v", err)
			return
		}

		if packet.Event != "chat:message" {
			t.Errorf("Sai tên sự kiện: mong muốn 'chat:message', nhận được '%s'", packet.Event)
		}

		var payload map[string]string
		if err := json.Unmarshal(packet.Data, &payload); err != nil {
			t.Errorf("Lỗi giải mã payload Data: %v", err)
			return
		}

		if payload["message"] != "hello gocket" {
			t.Errorf("Payload dữ liệu không trùng khớp: mong muốn 'hello gocket', nhận được '%s'", payload["message"])
		}
	}()

	// Gửi thử khi kết nối đang mở
	err := client.Emit("chat:message", testPayload)
	if err != nil {
		t.Fatalf("Emit thất bại khi kết nối đang hoạt động bình thường: %v", err)
	}

	// 2. Đóng kết nối và kiểm tra lỗi khi gọi Emit tiếp
	client.Close()

	err = client.Emit("chat:message", testPayload)
	if err == nil {
		t.Fatal("Lỗi logic: Emit đáng lẽ phải trả về lỗi sau khi client đã đóng kết nối")
	}
}
