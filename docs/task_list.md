# GOcket - Golang Websocket Realtime Engine

## Project Structure

```txt
Go-Socket-Clone/
├── frontend/             # NextJS, Tailwind, TypeScript (Client test)
├── backend/              # Golang, Gorilla WebSocket
├── docs/                 # Tài liệu, Architecture
├── docker-compose.yml    # Chạy Backend, Frontend
└── README.md
```

---

# DAY 1 — PROJECT SETUP & REST API BASE

- [ ] **TSK-001** `[PM/Setup]` Khởi tạo Monorepo Git + Base Structure. *(Estimate: 1h · Priority: Urgent)*

  **Description:**
  - Tạo repo Github.
  - Tạo 2 branch mặc định: `main`, `dev`.
  - Khởi tạo thư mục: `frontend/`, `backend/`, `docs/`.
  - Cập nhật file `.gitignore` cho Golang và Node.js.
  - *Lưu ý: Bắt đầu từ TSK-002, dev phải tạo nhánh `feature/TSK-xxx` từ nhánh `dev`.*

- [ ] **TSK-002** `[BE_Core]` Khởi tạo Golang module & REST Healthcheck. *(Estimate: 1h · Priority: Urgent)*

  **Description:**
  - Chạy `go mod init <tên-module>`.
  - Setup một HTTP Server cơ bản bằng `net/http`.
  - Tạo endpoint `GET /api/v1/health` để check server sống hay chết.
  - **Quy tắc (rules.md):** Response API phải bọc trong chuẩn:
    ```json
    {
      "success": true,
      "message": "Server is running",
      "data": null,
      "errorCode": null
    }
    ```

- [ ] **TSK-003** `[FE_Core]` Khởi tạo NextJS Client UI. *(Estimate: 1h · Priority: Urgent)*

  **Description:**
  - Chạy `npx create-next-app@latest frontend --typescript --tailwind --eslint --app`.
  - Dọn dẹp `page.tsx` và `globals.css`.
  - **Quy tắc (rules.md):** Setup Prettier & ESLint. Bắt buộc code theo chuẩn `camelCase` cho biến và `PascalCase` cho Component.

---

# DAY 2 — WEBSOCKET CORE & EVENT-BASED ARCHITECTURE

- [ ] **TSK-004** `[BE_WS]` Cài đặt Gorilla WebSocket & Upgrade HTTP. *(Estimate: 1.5h · Priority: Urgent)*

  **Description:**
  - Cài đặt `github.com/gorilla/websocket`.
  - Tạo endpoint `GET /ws`.
  - Khởi tạo struct `upgrader` với `CheckOrigin: return true` để Frontend NextJS gọi không bị CORS.
  - Validate lỗi upgrade: Nhớ quy tắc **phải check error rõ ràng `if err != nil`**.

- [ ] **TSK-005** `[BE_WS]` Định nghĩa Payload Models (Event-based). *(Estimate: 1h · Priority: Urgent)*

  **Description:**
  - Giả lập Socket.IO bằng cách bọc mọi tin nhắn vào struct JSON chung:
    ```go
    // Struct phải viết hoa (Exported) theo rules.md
    type WsMessage struct {
        Event string      `json:"event"`
        Data  interface{} `json:"data"` // Có thể map với bất kỳ JSON nào
    }
    ```

- [ ] **TSK-006** `[BE_WS]` Quản lý Client Memory & Global Broadcast. *(Estimate: 2.5h · Priority: Urgent)*

  **Description:**
  - Tạo biến nội bộ (unexported theo chuẩn `camelCase`): `var clients = make(map[*websocket.Conn]bool)`.
  - Tạo `broadcast = make(chan WsMessage)`.
  - Chạy 1 Goroutine (ví dụ: `go handleMessages()`) lắng nghe channel `broadcast` và đẩy message tới toàn bộ client (`WriteJSON`).
  - Xử lý lock (Mutex) nếu cần khi thao tác với map `clients` để tránh panic `concurrent map iteration`.

---

# DAY 3 — ADVANCED FEATURES (ROOMS & EVENT ROUTING)

- [ ] **TSK-007** `[BE_WS]` Viết logic Event Routing. *(Estimate: 2h · Priority: High)*

  **Description:**
  - Thay vì client gửi tin nhắn lên là đẩy hết vào `broadcast`. Hãy viết vòng lặp `switch msg.Event` (Event Router).
  - Ví dụ: 
    - `case "chat_message"`: Chạy hàm xử lý chat.
    - `case "ping"`: Trả về `pong`.
    - `default`: Bỏ qua hoặc log warning.

- [ ] **TSK-008** `[BE_WS]` Implement tính năng Rooms (Kênh riêng). *(Estimate: 3h · Priority: High)*

  **Description:**
  - Giống Socket.IO, client có thể `join` hoặc `leave` room.
  - Cấu trúc dữ liệu dự kiến: `var rooms = make(map[string]map[*websocket.Conn]bool)`.
  - Thêm logic xử lý event `join_room` (cập nhật client vào map room) và `leave_room`.

- [ ] **TSK-009** `[BE_WS]` Room Broadcasting. *(Estimate: 1.5h · Priority: High)*

  **Description:**
  - Thêm trường `RoomID` vào `WsMessage`.
  - Cập nhật hàm `handleMessages`: Nếu tin nhắn có `RoomID`, chỉ lặp qua map client của Room đó để `WriteJSON`, ngược lại thì Broadcast Global.

---

# DAY 4 — STABILITY, MEMORY LEAK & LIFECYCLE

- [ ] **TSK-010** `[BE_WS]` Xử lý Heartbeat (Ping/Pong). *(Estimate: 2h · Priority: Urgent)*

  **Description:**
  - WebSocket có thể bị treo connection (Client mất mạng đột ngột nhưng server không biết).
  - Dùng `SetReadDeadline`, `SetPongHandler` của Gorilla.
  - Từ server, định kỳ (vd 10s/lần) gửi một gói `Ping` xuống client. Nếu trong 15s không nhận được `Pong`, chủ động đóng connection.

- [ ] **TSK-011** `[BE_WS]` Clean up Memory & Disconnect Event. *(Estimate: 1.5h · Priority: Urgent)*

  **Description:**
  - Khi client ngắt kết nối (bắt được qua `ReadJSON` error hoặc Heartbeat timeout):
    - Đảm bảo xóa connection khỏi biến `clients` (global).
    - Duyệt qua biến `rooms` để xóa connection này khỏi mọi room mà nó đang join.
    - (Bắt buộc) Đóng `ws.Close()` để giải phóng socket descriptor.

---

# DAY 5 — FRONTEND INTEGRATION & UI TEST

- [ ] **TSK-012** `[FE_WS]` Xây dựng Custom Hook `useWebSocket`. *(Estimate: 2h · Priority: High)*

  **Description:**
  - Trong NextJS, tạo `hooks/useWebSocket.ts`.
  - Quản lý lifecycle của WebSocket client: Tự động connect khi component mount, tự động `ws.close()` khi unmount.
  - Tạo các hàm bọc (Wrapper functions) giống Socket.IO: `emit(event, data)` và `on(event, callback)`.

- [ ] **TSK-013** `[FE_UI]` Xây dựng Giao diện Chat Room. *(Estimate: 2h · Priority: Medium)*

  **Description:**
  - Code component `ChatPanel.tsx` (PascalCase theo `rules.md`).
  - Giao diện có:
    - Input nhập tên Room -> Nút "Join Room".
    - Khung hiển thị log chat.
    - Input tin nhắn + Nút "Gửi".

- [ ] **TSK-014** `[FE_WS]` Logic Tự động Reconnect (Auto-reconnect). *(Estimate: 1.5h · Priority: Medium)*

  **Description:**
  - Socket.IO có tính năng tự động nối lại rất ăn tiền.
  - Xử lý sự kiện `ws.onclose` ở Frontend: Nếu rớt mạng, chạy `setTimeout` thử kết nối lại mỗi 3 giây (Exponential backoff nếu cần).

---

# DAY 6 — DOCKERIZE, POLISH & DEPLOY

- [ ] **TSK-015** `[Deploy]` Dockerize Backend & Frontend. *(Estimate: 2h · Priority: Medium)*

  **Description:**
  - Viết `Dockerfile` cho Golang (Dùng base `golang:alpine`, build file binary).
  - Viết `Dockerfile` cho NextJS.
  - Viết `docker-compose.yml` để chạy cả 2 service chung 1 network. Expose port `8080` (Go) và `3000` (Next).

- [ ] **TSK-016** `[Docs]` Viết README.md & Hướng dẫn sử dụng. *(Estimate: 1h · Priority: Low)*

  **Description:**
  - Thêm hình ảnh/gif demo giao diện test chat real-time.
  - Document cách gửi JSON payload chuẩn để dev Frontend sau này biết cách xài:
    ```json
    { "event": "chat", "data": { "msg": "hello" } }
    ```

---

### 📝 Ghi chú tuân thủ Workflow dành cho Dev (Nhắc lại theo `first-step.md`)
1. Không code trực tiếp trên `main` hoặc `dev`.
2. Tạo nhánh: `git checkout -b feature/TSK-004`.
3. Khi commit phải theo chuẩn: `git commit -m "feat: [TSK-004] cài đặt gorilla websocket và upgrade HTTP"`.
4. Code xong push lên, tạo **Pull Request** vào nhánh `dev`, tag người khác review trước khi Merge.