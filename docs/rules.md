# QUY CHUẨN LÀM VIỆC CHUNG CỦA DỰ ÁN (PROJECT RULES & CONVENTIONS)

## 1. GIT WORKFLOW & BRANCHING
Dự án áp dụng Branching Model cơ bản của Github Flow:
*   **`main` branch:** Source code sạch tuyệt đối để Deploy lên Server Thật (Production). Chỉ Leader/PM mới được phép nhấn nút Merge vào nhánh này.
*   **`dev` branch:** Nhánh môi trường STAGING. Tất cả anh em dev sẽ ghép code vào đây để review chéo và test tổng thể trước khi release. Cấm push thẳng (direct push) vào nhánh này.
*   **Nhánh cá nhân (Feature / Fixbug):** Bắt buộc phải rẽ nhánh từ `dev`.
    *   Cú pháp nhánh: `<loại_nhánh>/TSK-<ID-task>`
    *   Ví dụ làm tính năng: `feature/TSK-010`
    *   Ví dụ sửa bug: `fixbug/TSK-010`
    *   Ví dụ viết docs/setup: `chore/TSK-001`

## 2. QUY TẮC COMMIT MESSAGE (CONVENTIONAL COMMITS)
Bắt buộc phải có Tiền tố và Mã Task để dễ theo dõi lịch sử:
*   `feat: [TSK-010] tạo endpoint GET /api/v1/route` *(Thêm tính năng mới)*
*   `fix: [TSK-020] sửa lỗi hiển thị sai icon marker` *(Sửa lỗi)*
*   `chore: [TSK-002] cập nhật NextJS và cài Tailwind` *(Cấu hình, thư viện, không tác động code logic)*
*   `refactor: [TSK-015] tối ưu lại vòng lặp Dijkstra` *(Sửa code nhưng không làm thay đổi tính năng)*
*   `docs: [TSK-029] thêm file hướng dẫn setup Docker` *(Cập nhật tài liệu)*

---

## 3. GIT STEP-BY-STEP (QUY TRÌNH HÀNG NGÀY)

Quy trình chuẩn khi bắt tay vào làm một Task mới:

```bash
# 1. Chuyển về nhánh dev và cập nhật code mới nhất từ team
git checkout dev
git pull origin dev

# 2. Tạo nhánh mới cho task của mình
git checkout -b feature/TSK-010


# 3. Add và Commit code
git add .
git commit -m "feat: [TSK-010] implement shortest path API"

# 4. Push nhánh cá nhân lên Github
git push -u origin feature/TSK-010
```
⚠️ **Quan trọng:** Sau khi push, lên Github tạo một Pull Request (PR) từ nhánh feature/TSK-010 vào nhánh dev. Gắn thẻ (Tag/Assign) một thành viên khác trong team để Review Code. Review xong mới được bấm Merge.

## 4. REST API STANDARDS (BACKEND)

### Cấu trúc JSON Response (Chuẩn 1 chiều)
  - Tất cả API response (dù thành công hay thất bại) phải được bọc vào một Interface/DTO duy nhất để Frontend dễ dàng parse JSON.
```json
{
  "success": true, // true | false
  "message": "Lấy lộ trình thành công", 
  "data": {       // Payload trả về (nếu mảng trống trả [], nếu không có data trả null)
     "distance": 1200.5,
     "path": [[10.776, 106.700], [10.778, 106.702]]
  },
  "errorCode": null // Mã lỗi nội bộ để Frontend map UI (VD: "ERR_OUT_OF_BOUNDS"), không có lỗi thì null
}
```

### Quy tắc định tuyến (Routing)
  - Dùng danh từ số nhiều, viết thường (lowercase), phân cách bằng dấu gạch
    ngang (kebab-case).
  - **Đúng:** `GET /api/v1/routes`, `GET /api/v1/search-locations`
  - **Sai:** `GET /api/v1/getRoute`, `GET /api/v1/Search`

### Quy tắc HTTP Status Code
  - `200 OK`: Trả về thành công (Dùng cho GET, PUT, DELETE).
  - `201 Created`: Tạo mới thành công (Dùng cho POST).
  - `400 Bad Request`: Client gửi sai data, thiếu query params.
  - `404 Not Found`: Không tìm thấy dữ liệu (VD: Không tìm thấy đường đi).
  - `500 Internal Server Error`: Lỗi server (DB sập, code Golang panic...).

## 5. CODE STYLE & CONVENTIONS

### Đối với Frontend (NextJS/TypeScript)
  - **Biến và Hàm:** Dùng `camelCase` (vd: `findShortestRoute`, `startPoint`).
  - **Component & File (React):** Dùng `PascalCase` (vd: `MapView.tsx`, `SearchPanel.tsx`).
  - **Type / Interface:** Dùng `PascalCase`, ưu tiên dùng `interface` thay cho `type`.
    (vd: interface RouteResponse {}).
  - **Bắt buộc** chạy ESLint & Prettier trước khi commit code. Không được để sót console.log() lên môi trường dev.

### Đối với Backend (Golang)
  - Bắt buộc chạy lệnh `go fmt ./...` để format code trước khi commit.
  - **Struct / Function public (Exported):** Bắt buộc viết hoa chữ cái đầu PascalCase (vd: `type Node struct`, `func BuildGraph()`).
  - **Biến nội bộ (Unexported):** Dùng `camelCase` (vd: `routeCache`).
  - **Xử lý lỗi (Error Handling):** Phải check error rõ ràng `if err != nil`, không dùng `_` để bỏ qua lỗi trừ khi thật sự chắc chắn.

### Đối với Database (PostgreSQL / PostGIS)
  - **Tên Bảng (Table) và Cột (Column):** Dùng `snake_case` chữ thường (vd: `planet_osm_line`, `start_node_id`).
  - Không dùng tiếng Việt có dấu, không dùng khoảng trắng trong Database.

---