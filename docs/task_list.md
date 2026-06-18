# Go-Socket — MVP — Task List

> **Mô tả dự án:** Xây dựng một thư viện WebSocket bằng **Golang**, gồm các cơ chế cốt lõi (event-based messaging, room, namespace, broadcast, heartbeat, acknowledgement).

---

## Project Structure

```txt
go-socket/
├── server/
│   ├── client.go
│   ├── hub.go
│   ├── room.go
│   ├── namespace.go
│   └── event.go
├── examples/
│   └── chat-demo/
├── docs/
├── tests/
├── go.mod
└── README.md
```

> Convention tuân theo `rules.md`:
> - Branch: `feature/TSK-xxx` · `fixbug/TSK-xxx` · `chore/TSK-xxx`
> - Commit: `feat:` `fix:` `chore:` `refactor:` `docs:` kèm `[TSK-xxx]`
> - Go: `go fmt ./...` trước commit · Exported → `PascalCase` · internal → `camelCase` · luôn check `if err != nil`

---

# PHASE 1 — PROJECT SETUP + WEBSOCKET CƠ BẢN

- [ ] **TSK-001** `[PM/Setup]` Khởi tạo Go Module + Base Structure. *(Estimate: 1h · Priority: Urgent)*

  **Description:**
  - Tạo repo Github `gocket`, tạo branch `main`, `dev`
  - Chạy `go mod init github.com/apase95/gocket.git`
  - Tạo folder: `server/`, `examples/`, `docs/`, `tests/`
  - Tạo README mô tả: Mục tiêu, Stack, Roadmap

- [ ] **TSK-002** `[Core]` WebSocket Upgrade (HTTP → WS). *(Estimate: 30m · Priority: Urgent)*

  **Description:**
  - Cài `gorilla/websocket`: `go get github.com/gorilla/websocket`
  - Viết hàm `Upgrade(w http.ResponseWriter, r *http.Request) (*websocket.Conn, error)` trong `server/hub.go`
  - Cấu hình `websocket.Upgrader` (ReadBufferSize, WriteBufferSize, CheckOrigin)

- [ ] **TSK-003** `[Core]` Định nghĩa struct `Client`. *(Estimate: 1.5h · Priority: Urgent)*

  **Description:**
  - File `server/client.go`:
    ```go
    type Client struct {
        ID          string
        Conn        *websocket.Conn
        Send        chan []byte
        hub         *Hub
        rooms       map[string]bool
        connectedAt time.Time
    }
    ```
  - Sinh `ID` tự động bằng `uuid`: `go get github.com/google/uuid`

- [ ] **TSK-004** `[Core]` Định nghĩa struct `Hub`. *(Estimate: 2h · Priority: Urgent)*

  **Description:**
  - File `server/hub.go`:
    ```go
    type Hub struct {
        clients    map[string]*Client
        register   chan *Client
        unregister chan *Client
        mu         sync.RWMutex
    }
    ```
  - Implement `Run()` chạy goroutine loop xử lý `register`/`unregister`
  - Mọi thao tác đọc/ghi `clients` map phải dùng `sync.RWMutex`

- [ ] **TSK-005** `[Core]` Read Pump. *(Estimate: 1.5h · Priority: Urgent)*

  **Description:**
  - `func (c *Client) readPump()` chạy trong goroutine riêng
  - Loop `conn.ReadMessage()` → đẩy raw bytes vào Dispatcher
  - Khi đọc lỗi (disconnect, timeout): `defer hub.unregister <- client` + `conn.Close()`

- [ ] **TSK-006** `[Core]` Write Pump. *(Estimate: 1.5h · Priority: Urgent)*

  **Description:**
  - `func (c *Client) writePump()` chạy trong goroutine riêng
  - Loop chờ `c.Send` channel → `conn.WriteMessage()`
  - Nếu `c.Send` channel đóng → gửi `CloseMessage` rồi return
  - Dùng `sync.Once` để đảm bảo channel chỉ đóng 1 lần (tránh panic)

---

# PHASE 2 — EVENT ENGINE

- [ ] **TSK-007** `[Core]` Định nghĩa `Packet` (Event Payload). *(Estimate: 1h · Priority: Urgent)*

  **Description:**
  - File `server/event.go`:
    ```go
    type Packet struct {
        Event string          `json:"event"`
        Data  json.RawMessage `json:"data,omitempty"`
        AckID string          `json:"ackId,omitempty"`
    }
    ```
  - Mọi message qua WS đều encode/decode theo format `Packet` này

- [ ] **TSK-008** `[Core]` Implement `Emit()`. *(Estimate: 1h · Priority: Urgent)*

  **Description:**
  - `func (c *Client) Emit(event string, data any) error`
  - Marshal `data` → `Packet` → JSON → đẩy vào `c.Send`
  - Trả lỗi nếu client đã ngắt kết nối

- [ ] **TSK-009** `[Core]` Implement `On()` — Đăng ký Event Handler. *(Estimate: 1.5h · Priority: Urgent)*

  **Description:**
  - Hub giữ `handlers map[string]HandlerFunc`
  - `type HandlerFunc func(c *Client, data json.RawMessage)`
  - `func (h *Hub) On(event string, handler HandlerFunc)`

- [ ] **TSK-010** `[Core]` Implement Dispatcher. *(Estimate: 2h · Priority: Urgent)*

  **Description:**
  - Trong `readPump()`: decode raw bytes → `Packet`
  - Lookup `hub.handlers[packet.Event]`, nếu có thì gọi handler
  - Nếu event không có handler: `log.Printf("unknown event: %s", packet.Event)`, không crash
  - Nếu JSON decode lỗi: gửi `Packet{Event: "error"}` về cho client

---

# PHASE 3 — ROOM

- [ ] **TSK-011** `[Core]` Định nghĩa `Room` + Join/Leave/LeaveAll. *(Estimate: 2h · Priority: Urgent)*

  **Description:**
  - File `server/room.go`:
    ```go
    type Room struct {
        Name    string
        clients map[string]*Client
        mu      sync.RWMutex
    }
    ```
  - Hub quản lý `rooms map[string]*Room` với `sync.RWMutex`
  - `Join(client, roomName)`: tạo room nếu chưa tồn tại, thêm client vào
  - `Leave(client, roomName)`: xóa client, nếu room rỗng → xóa room khỏi Hub (tránh memory leak)
  - `LeaveAll(client)`: gọi tự động khi client disconnect

- [ ] **TSK-012** `[Core]` Broadcast tới Room. *(Estimate: 1.5h · Priority: Urgent)*

  **Description:**
  - `func (h *Hub) BroadcastToRoom(roomName, event string, data any)`
  - Loop qua `room.clients`, gọi `client.Emit()` — bỏ qua lỗi từng client
  - `BroadcastExcept(roomName, excludeID, event, data)` để skip người gửi

---

# PHASE 4 — NAMESPACE + BROADCAST

- [ ] **TSK-013** `[Core]` Implement Namespace. *(Estimate: 2h · Priority: High)*

  **Description:**
  - File `server/namespace.go`:
    ```go
    type Namespace struct {
        Path string
        hub  *Hub
    }
    ```
  - Mỗi namespace có Hub riêng → client, room, handler không chia sẻ với namespace khác
  - Đăng ký route: `http.HandleFunc("/ws/chat", ns.HandleWS)`

- [ ] **TSK-014** `[Core]` `BroadcastAll()` + `EmitTo()`. *(Estimate: 1.5h · Priority: Urgent)*

  **Description:**
  - `func (h *Hub) BroadcastAll(event string, data any)`: gửi tới toàn bộ client trong Hub
    - Dùng `sync.WaitGroup` + goroutine để không blocking
  - `func (h *Hub) EmitTo(clientID, event string, data any) error`: gửi tới đúng 1 client
    - Trả `ErrClientNotFound` nếu ID không tồn tại

---

# PHASE 5 — HEARTBEAT + ACK

- [ ] **TSK-015** `[Core]` Heartbeat (Ping/Pong + Auto Disconnect). *(Estimate: 2h · Priority: Urgent)*

  **Description:**
  - `writePump()` dùng `time.NewTicker(30s)` để gửi `PingMessage` định kỳ
  - `conn.SetPongHandler()`: reset `ReadDeadline` mỗi khi nhận Pong
  - `conn.SetReadDeadline(time.Now().Add(60s))`: nếu không nhận Pong trong 60s → deadline expire → `readPump` trả lỗi → client bị unregister

- [ ] **TSK-016** `[Core]` Acknowledgement (Ack). *(Estimate: 2.5h · Priority: High)*

  **Description:**
  - Mỗi Client có `pendingAcks map[string]chan json.RawMessage`
  - Dispatcher: nếu `packet.Event == "ack"` → resolve channel theo `packet.AckID`
  - `func (c *Client) EmitWithAck(event string, data any, timeout time.Duration) (json.RawMessage, error)`
    - Tạo `AckID`, emit, dùng `select` chờ channel hoặc `time.After(timeout)`
    - Khi client disconnect → cleanup toàn bộ pending ack (tránh memory leak)

---

# PHASE 6 — DEMO + TEST + DOCS

- [ ] **TSK-017** `[Demo]` Chat Demo (HTML/JS Client). *(Estimate: 2h · Priority: High)*

  **Description:**
  - `examples/chat-demo/`: server Go + `index.html` dùng WebSocket API thuần (browser)
  - Connect tới `ws://localhost:8080/ws/chat`
  - UI tối thiểu: chọn phòng, input gửi message, list message realtime, list user online

- [ ] **TSK-018** `[Testing]` Unit Test. *(Estimate: 2h · Priority: High)*

  **Description:**
  - Test `Hub`: register/unregister client, map `clients` cập nhật đúng
  - Test `Room`: Join/Leave, room tự xóa khi rỗng
  - Test `Dispatcher`: event đã đăng ký → handler gọi đúng 1 lần; event lạ → không crash
  - Test `Ack`: resolve đúng, trả lỗi timeout đúng
  - Bắt buộc chạy với flag: `go test ./... -race`

- [ ] **TSK-019** `[Docs]` README + Packet Convention. *(Estimate: 1.5h · Priority: Medium)*

  **Description:**
  - Quickstart: cài đặt, khởi tạo Hub, `On`, `Emit`, `Join`, chạy server
  - Sơ đồ kiến trúc: `Client ↔ WS ↔ Hub ↔ Dispatcher ↔ Room/Namespace`
  - Document `Packet` JSON format
  - Naming convention event: `domain:action` (vd: `chat:message`, `user:join`)
  - Liệt kê built-in events: `connect`, `disconnect`, `error`, `ack`, `server:closing`

---

# MVP Done Checklist

- [ ] Server WS chạy được, accept connection qua `ws://`
- [ ] `Emit()` / `On()` hoạt động đúng — gửi/nhận event 2 chiều
- [ ] Join/Leave Room + Broadcast đúng phạm vi (room / all)
- [ ] Namespace: `/ws/chat` và `/ws/admin` độc lập nhau
- [ ] Heartbeat tự động disconnect client "chết"
- [ ] `EmitWithAck()` hoạt động với timeout, không memory leak
- [ ] Chat demo chạy thực tế nhiều tab/nhiều client
- [ ] `go test ./... -race` pass sạch — không race condition, không goroutine leak

---

# Post-MVP (Optional)

- [ ] Middleware (pre-connection) — kiểm tra token trước khi upgrade WS
- [ ] Graceful Shutdown — broadcast `server:closing`, chờ client ngắt trước khi tắt
- [ ] Structured Logging (`log/slog`) thay `log.Printf`
- [ ] Management API — `GET /api/v1/clients`, `GET /api/v1/rooms`
- [ ] Load Test — 500+ connection đồng thời, đo broadcast latency
- [ ] Docker + docker-compose
- [ ] Redis Adapter — scale-out nhiều instance (Pub/Sub cross-node)
- [ ] JWT Authentication Middleware
- [ ] Rate Limiting — giới hạn messages/giây per client
- [ ] Auto-Reconnect Client SDK (JS/TS) — exponential backoff
- [ ] Metrics — expose `/debug/vars` (active connections, messages sent/received)

---
