# Go-Socket-Lite MVP — Task List

> **Mô tả dự án:** Xây dựng một thư viện WebSocket bằng **Golang**, gồm các cơ chế cốt lõi (event-based messaging, room, namespace, broadcast, heartbeat, acknowledgement).

---

## Project Structure

```txt
go-socket-lite/
├── server/               # Core: Hub, Client, Room, Namespace, Event Engine
│   ├── hub.go
│   ├── client.go
│   ├── room.go
│   ├── namespace.go
│   ├── event.go
│   └── middleware.go
├── examples/
│   └── chat-demo/        # Demo HTML/JS client + main.go chạy thử
├── docs/                 # Tài liệu, kiến trúc, hướng dẫn sử dụng
├── tests/                # Unit test + load test script
├── go.mod
├── docker-compose.yml
└── README.md
```

---

# DAY 1 — PROJECT SETUP + WEBSOCKET CƠ BẢN

- [ ] **TSK-001** `[PM/Setup]` Khởi tạo Go Module + Base Structure. *(Estimate: 1h · Priority: Urgent)*

  **Description:**
  - Tạo repo Github `go-socket-lite`
  - Tạo branch: `main`, `dev`
  - Chạy `go mod init github.com/<org>/go-socket-lite`
  - Tạo folder: `server/`, `examples/`, `docs/`, `tests/`
  - Tạo README mô tả: Mục tiêu, Stack, Roadmap

- [ ] **TSK-002** `[Core]` Cài đặt thư viện WebSocket cơ bản. *(Estimate: 30m · Priority: Urgent)*

  **Description:**
  - Cài `gorilla/websocket`:
    ```bash
    go get github.com/gorilla/websocket
    ```
  - Viết hàm `Upgrade(w http.ResponseWriter, r *http.Request) (*websocket.Conn, error)` để upgrade HTTP → WS.

- [ ] **TSK-003** `[Core]` Định nghĩa struct `Client`. *(Estimate: 1.5h · Priority: Urgent)*

  **Description:**
  - File `server/client.go`:
    ```go
    type Client struct {
        ID       string
        Conn     *websocket.Conn
        Send     chan []byte
        Hub      *Hub
        Rooms    map[string]bool
    }
    ```
  - Sinh `ID` tự động (uuid) khi client connect.

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
  - Implement `Run()` chạy goroutine loop xử lý `register` / `unregister`.
  - Đảm bảo thread-safe khi đọc/ghi map `clients` (dùng `sync.RWMutex`).

- [ ] **TSK-005** `[Core]` Implement Read/Write Pump cho `Client`. *(Estimate: 2h · Priority: Urgent)*

  **Description:**
  - `readPump()`: loop đọc message từ `conn.ReadMessage()`, đẩy vào dispatcher.
  - `writePump()`: loop đọc từ channel `Send`, ghi ra `conn.WriteMessage()`.
  - Xử lý `defer` đóng connection + gọi `hub.unregister <- client` khi client disconnect.
  - Mỗi pump chạy trong goroutine riêng (`go client.readPump()`, `go client.writePump()`).

---

# DAY 2 — EVENT ENGINE (EMIT / ON)

- [ ] **TSK-006** `[Core]` Định nghĩa `Packet` (Event Payload). *(Estimate: 1h · Priority: Urgent)*

  **Description:**
  - File `server/event.go`:
    ```go
    type Packet struct {
        Event string          `json:"event"`
        Data  json.RawMessage `json:"data"`
        AckID string          `json:"ackId,omitempty"`
    }
    ```
  - Mọi message qua WS đều phải encode/decode theo format `Packet` này (giống "envelope" trong `rules.md`).

- [ ] **TSK-007** `[Core]` Implement `Emit()` gửi event tới Client. *(Estimate: 1.5h · Priority: Urgent)*

  **Description:**
  - Hàm `func (c *Client) Emit(event string, data interface{}) error`
  - Marshal `data` → `Packet` → JSON → đẩy vào `c.Send`.
  - Trả lỗi nếu `c.Send` đầy (channel buffer full) hoặc client đã disconnect.

- [ ] **TSK-008** `[Core]` Implement `On()` đăng ký Event Handler. *(Estimate: 2h · Priority: Urgent)*

  **Description:**
  - Hub giữ map `handlers map[string]func(*Client, json.RawMessage)`.
  - Hàm `func (h *Hub) On(event string, handler func(*Client, json.RawMessage))`.
  - Ví dụ sử dụng:
    ```go
    hub.On("chat:message", func(c *Client, data json.RawMessage) {
        hub.BroadcastAll("chat:message", data)
    })
    ```

- [ ] **TSK-009** `[Core]` Implement Dispatcher. *(Estimate: 2h · Priority: Urgent)*

  **Description:**
  - Trong `readPump()`, decode raw bytes → `Packet`.
  - Lookup `hub.handlers[packet.Event]`, nếu có thì gọi handler.
  - Nếu event không tồn tại → log warning, không crash server.
  - Validate input: nếu JSON decode lỗi → trả về `Packet{Event: "error", Data: ...}` cho client.

- [ ] **TSK-010** `[Core]` Middleware Layer (Pre-connection). *(Estimate: 1.5h · Priority: Medium)*

  **Description:**
  - Định nghĩa `type Middleware func(r *http.Request) (clientMeta map[string]string, err error)`.
  - Hub có `Use(mw Middleware)`, chạy trước khi `Upgrade()`.
  - Use case: kiểm tra token trong query param `?token=xxx`, nếu lỗi → trả `401`, không upgrade connection.

---

# DAY 3 — ROOMS & NAMESPACES

- [ ] **TSK-011** `[Core]` Định nghĩa `Room` + Join/Leave. *(Estimate: 2h · Priority: Urgent)*

  **Description:**
  - File `server/room.go`:
    ```go
    type Room struct {
        Name    string
        Clients map[string]*Client
        mu      sync.RWMutex
    }
    ```
  - Hub quản lý `rooms map[string]*Room`.
  - Hàm `func (h *Hub) Join(client *Client, roomName string)` và `Leave(...)`.
  - Khi client disconnect → tự động leave tất cả room đang tham gia.

- [ ] **TSK-012** `[Core]` Broadcast tới Room. *(Estimate: 1.5h · Priority: Urgent)*

  **Description:**
  - Hàm `func (h *Hub) BroadcastToRoom(roomName, event string, data interface{})`.
  - Loop qua `room.Clients`, gọi `client.Emit()` cho từng client (skip lỗi, không để 1 client lỗi làm crash broadcast).
  - Option `Except(clientID string)` để broadcast trừ chính người gửi.

- [ ] **TSK-013** `[Core]` Implement Namespace. *(Estimate: 2.5h · Priority: High)*

  **Description:**
  - File `server/namespace.go`. Mỗi `Namespace` là 1 Hub con, đăng ký theo path: `/ws/chat`, `/ws/admin`.
  - `type Namespace struct { Path string; Hub *Hub }`
  - Server đăng ký route `http.HandleFunc("/ws/chat", ns.HandleWS)`.
  - Mỗi namespace có handler/event riêng, không chia sẻ state với namespace khác.

- [ ] **TSK-014** `[Core]` Client Tracking (Room/Namespace mapping). *(Estimate: 1h · Priority: Medium)*

  **Description:**
  - Thêm field `Rooms map[string]bool` và `Namespace string` vào `Client`.
  - Viết hàm tiện ích `func (c *Client) Rooms() []string` trả về danh sách room hiện tại (debug/API).

---

# DAY 4 — HEARTBEAT, ACK, GRACEFUL SHUTDOWN

- [ ] **TSK-015** `[Core]` Heartbeat (Ping/Pong) + Timeout Disconnect. *(Estimate: 2h · Priority: Urgent)*

  **Description:**
  - Set `conn.SetReadDeadline()` + `conn.SetPongHandler()`.
  - `writePump()` gửi `PingMessage` định kỳ (vd: 30s) bằng `time.Ticker`.
  - Nếu không nhận `Pong` trong thời gian timeout (vd: 60s) → đóng connection, `unregister` client.

- [ ] **TSK-016** `[Core]` Implement Acknowledgement (Ack Callback). *(Estimate: 2.5h · Priority: High)*

  **Description:**
  - Khi gửi event kèm `AckID`, lưu callback vào map `pendingAcks map[string]chan json.RawMessage`.
  - Khi client gửi packet với `event = "ack"` và `ackId` trùng → resolve channel tương ứng.
  - Hàm `func (c *Client) EmitWithAck(event string, data interface{}, timeout time.Duration) (json.RawMessage, error)`.

- [ ] **TSK-017** `[Core]` Graceful Shutdown cho Hub + Client. *(Estimate: 1.5h · Priority: Medium)*

  **Description:**
  - Implement `func (h *Hub) Shutdown(ctx context.Context) error`.
  - Khi nhận `SIGINT`/`SIGTERM`: gửi event `server:closing` tới tất cả client, đóng connection sau X giây.
  - Đảm bảo không leak goroutine (mọi pump phải return khi channel đóng).

- [ ] **TSK-018** `[Core]` Error Handling + Logging chuẩn. *(Estimate: 1h · Priority: Medium)*

  **Description:**
  - Tạo package `logger` (dùng `log/slog` của Go standard library).
  - Mọi error trong `readPump`/`writePump`/`dispatcher` phải log kèm `client.ID`.
  - Tuân thủ `if err != nil { ... }` rõ ràng, không bỏ qua lỗi bằng `_`.

---

# DAY 5 — BROADCAST APIs & MANAGEMENT REST API

- [ ] **TSK-019** `[Core]` `BroadcastAll()` — Gửi tới toàn bộ Client. *(Estimate: 1h · Priority: Urgent)*

  **Description:**
  - `func (h *Hub) BroadcastAll(event string, data interface{})`.
  - Loop toàn bộ `h.clients`, `Emit()` cho từng client, dùng `goroutine` + `sync.WaitGroup` nếu số client lớn.

- [ ] **TSK-020** `[Core]` `EmitTo(clientID, event, data)`. *(Estimate: 1h · Priority: High)*

  **Description:**
  - Gửi event tới đúng 1 client theo `ID`.
  - Trả lỗi `ErrClientNotFound` nếu `clientID` không tồn tại trong `h.clients`.

- [ ] **TSK-021** `[BE_API]` Management API — `GET /api/v1/clients`. *(Estimate: 1.5h · Priority: Medium)*

  **Description:**
  - Theo chuẩn response của `rules.md` (`success`, `message`, `data`, `errorCode`).
  - Trả về danh sách client đang connect: `id`, `namespace`, `rooms`, `connectedAt`.

- [ ] **TSK-022** `[BE_API]` Management API — `GET /api/v1/rooms`. *(Estimate: 1h · Priority: Medium)*

  **Description:**
  - Trả về danh sách room hiện có + số lượng client trong từng room.
  - Routing tuân thủ kebab-case + danh từ số nhiều: `/api/v1/rooms`, `/api/v1/clients`.

---

# DAY 6 — DEMO CLIENT + TESTING

- [ ] **TSK-023** `[Demo]` Viết Demo Chat App (HTML/JS Client). *(Estimate: 2h · Priority: High)*

  **Description:**
  - Folder `examples/chat-demo/`.
  - File `index.html` dùng `WebSocket` API thuần (browser) để connect `ws://localhost:8080/ws/chat`.
  - UI tối thiểu: input gửi message, list hiển thị message realtime, list user trong room.

- [ ] **TSK-024** `[Testing]` Unit Test cho `Hub`. *(Estimate: 2h · Priority: High)*

  **Description:**
  - Test `register`/`unregister` client (đảm bảo map `clients` cập nhật đúng).
  - Test `Join`/`Leave` room (đảm bảo client xuất hiện/biến mất khỏi `room.Clients`).
  - Dùng `go test ./server/... -race` để detect race condition.

- [ ] **TSK-025** `[Testing]` Unit Test cho Event Dispatcher + Ack. *(Estimate: 1.5h · Priority: Medium)*

  **Description:**
  - Test: gửi packet với event đã đăng ký → handler được gọi đúng 1 lần.
  - Test: gửi packet với event không tồn tại → không crash, trả packet `error`.
  - Test `EmitWithAck()` với timeout (đảm bảo trả lỗi đúng khi không có ack).

- [ ] **TSK-026** `[Testing]` Load Test — Nhiều Connection Đồng Thời. *(Estimate: 1.5h · Priority: Medium)*

  **Description:**
  - Viết script Go (`tests/loadtest/main.go`) mở N connection đồng thời (vd: 500-1000), mỗi connection gửi message định kỳ.
  - Đo: thời gian broadcast tới toàn bộ client, memory/goroutine usage (`pprof`).
  - Mục tiêu: server xử lý ổn định ≥ 500 connection trên máy local.

---

# DAY 7 — POLISH + DOCKER + DOCS

- [ ] **TSK-027** `[Deploy]` Dockerize Server. *(Estimate: 1.5h · Priority: Medium)*

  **Description:**
  - Viết `Dockerfile` multi-stage build (builder Go → image alpine nhỏ gọn).
  - Expose port `8080`.

- [ ] **TSK-028** `[Deploy]` Hoàn thiện `docker-compose.yml`. *(Estimate: 1h · Priority: Medium)*

  **Description:**
  - Service `go-socket-lite` chạy server chính.
  - (Optional) Service `redis` nếu chuẩn bị cho Post-MVP (Adapter scale-out).
  - Chạy `docker compose up` lên toàn bộ stack thành công.

- [ ] **TSK-029** `[Docs]` Viết README + Hướng dẫn sử dụng API. *(Estimate: 1.5h · Priority: Medium)*

  **Description:**
  - Hướng dẫn: cài đặt, khởi tạo `Hub`, `Emit`, `On`, `Join/Leave Room`, `Namespace`.
  - Vẽ sơ đồ kiến trúc đơn giản (Client ↔ WS ↔ Hub ↔ Rooms/Namespaces).
  - Ví dụ code mẫu cho cả Server (Go) và Client (JS).

- [ ] **TSK-030** `[Docs]` Tài liệu Packet Format & Event Convention. *(Estimate: 1h · Priority: Low)*

  **Description:**
  - Document chuẩn `Packet` JSON (giống chuẩn response trong `rules.md`).
  - Định nghĩa naming convention cho event: dùng `snake:case` hoặc `domain:action` (vd: `chat:message`, `user:join`, `room:leave`).
  - Liệt kê các event hệ thống built-in: `connect`, `disconnect`, `error`, `server:closing`.

---

# MVP Done Checklist

- [ ] Server WS chạy được, accept connection qua `ws://`.
- [ ] `Emit()` / `On()` hoạt động đúng — gửi/nhận event 2 chiều.
- [ ] Join/Leave Room hoạt động, Broadcast đúng phạm vi (room/namespace/all).
- [ ] Heartbeat tự động disconnect client "chết" (mất kết nối không graceful).
- [ ] Acknowledgement (Ack) hoạt động với timeout.
- [ ] Demo Chat App chạy thực tế nhiều tab/nhiều client.
- [ ] Toàn bộ test (`go test ./... -race`) pass, không leak goroutine.
- [ ] `docker compose up` chạy server thành công.

---

# Post-MVP (Optional)

- [ ] Redis Adapter — cho phép scale-out nhiều instance server (Pub/Sub cross-node broadcast).
- [ ] Binary Data Support — gửi/nhận `[]byte` (ngoài JSON) qua event.
- [ ] JWT Authentication Middleware — xác thực token trước khi upgrade WS.
- [ ] Auto-Reconnect Client SDK (JS/TS) — tự reconnect kèm exponential backoff.
- [ ] Rate Limiting — giới hạn số message/giây mỗi client để chống spam/abuse.
- [ ] Metrics/Monitoring — expose `/metrics` (Prometheus) cho số connection, throughput.
- [ ] Volatile Events — event không cần đảm bảo delivery (giống `socket.volatile.emit`).

---
