# HƯỚNG DẪN BƯỚC ĐẦU (FIRST STEP)

---

## BƯỚC 0: CHUẨN BỊ
1. Đảm bảo máy bạn đã cài đặt [Git](https://git-scm.com/).
2. **[QUAN TRỌNG]** Nhắn username Github của bạn để add quyền **Collaborator**.
3. Mở Terminal (hoặc Git Bash / VS Code Terminal) lên để bắt đầu.

---

## BƯỚC 1: CLONE DỰ ÁN VỀ MÁY
Để lấy toàn bộ source code từ Github về máy tính của bạn, dùng lệnh `clone`.

```bash
# Clone dự án về máy
git clone https://github.com/apase95/gocket.git

# Di chuyển vào thư mục dự án vừa tải về
cd gocket
```

*Lưu ý: Khi clone, Git đã tự động thiết lập kết nối (remote) tới Github với tên mặc định là `origin`.*

---

## BƯỚC 2: CHUYỂN SANG NHÁNH `dev` VÀ CẬP NHẬT CODE
Nhánh `dev` là nơi chứa code mới nhất của cả team. Mặc định khi clone về, bạn đang ở nhánh `main`. Hãy chuyển sang `dev` và cập nhật.

```bash
# Chuyển sang nhánh dev
git checkout dev

# Kéo (pull) code mới nhất từ kho lưu trữ (origin) nhánh dev về máy
git pull origin dev
```

---

## BƯỚC 3: TẠO NHÁNH LÀM VIỆC CÁ NHÂN
**⚠️ LUẬT CỦA TEAM:** Tuyệt đối KHÔNG code trực tiếp trên nhánh `main` hoặc `dev`. Bạn phải tạo một nhánh riêng từ `dev` để làm task của mình.

Giả sử bạn được giao task số 10 (TSK-010), hãy tạo nhánh mới:

```bash
# Lệnh -b giúp tạo nhánh mới VÀ chuyển sang nhánh đó luôn
git checkout -b feature/TSK-010
```
*(Nếu bạn fix bug, hãy đặt tên là `fixbug/TSK-xxx`)*

---

## BƯỚC 4: CODE, ADD VÀ COMMIT
Bây giờ bạn bắt đầu mở code lên và làm task của mình (tạo file mới, sửa code...). 
Sau khi code xong và chạy thử ngon lành:

```bash
# 1. Kiểm tra xem mình đã sửa những file nào
git status

# 2. Đưa TẤT CẢ các file đã sửa vào trạng thái chờ (Staging)
git add .
# Hoặc nếu chỉ muốn add từng file: git add frontend/package.json

# 3. Đóng gói code (Commit) kèm theo lời nhắn theo chuẩn của team
git commit -m "feat: [TSK-010] tạo giao diện tìm kiếm đường đi"
```
*(Nhớ tuân thủ quy tắc ghi chú Commit: dùng `feat:`, `fix:`, `chore:`, `docs:` kèm mã Task).*

---

## BƯỚC 5: PUSH CODE LÊN GITHUB (LẦN ĐẦU TIÊN)
Vì nhánh `feature/TSK-010` chỉ mới tồn tại trên máy tính của bạn, Github chưa hề biết đến nó. Lần đầu tiên đẩy code lên, bạn phải dùng cờ `-u` (upstream) để liên kết nhánh ở máy tính với nhánh trên Github.

```bash
git push -u origin feature/TSK-010
```
*Từ những lần push sau trên cùng nhánh này, bạn chỉ cần gõ ngắn gọn: `git push`.*

---

## BƯỚC 6: TẠO PULL REQUEST (PR) ĐỂ GỘP CODE
Code của bạn đã lên Github, nhưng nó vẫn nằm ở nhánh riêng của bạn. Để đưa code vào nhánh chung `dev`:

1. Lên trang Github của dự án.
2. Bạn sẽ thấy một nút màu xanh lá nổi bật: **"Compare & pull request"**. Bấm vào đó.
3. Đảm bảo nhánh gốc (base) là `dev`, nhánh so sánh (compare) là nhánh của bạn `feature/TSK-010`.
4. Viết mô tả ngắn gọn những gì bạn đã làm.
5. Ở góc phải, mục **Reviewers**, hãy tag (chọn) tên một người bạn trong team để họ xem code giúp bạn.
6. Bấm **Create pull request**.

---

## MỘT SỐ QUY TẮC PHẢI NHỚ

*   **Luôn đồng bộ trước khi push:** Nếu task của bạn làm trong nhiều ngày, nhánh `dev` có thể đã được người khác cập nhật code mới. Thỉnh thoảng hãy chạy lệnh sau để kéo code mới từ `dev` vào nhánh của bạn, tránh bị conflict (xung đột):
    ```bash
    git pull origin dev
    ```
*   **Xoá nhánh sau khi xong việc:** Khi Pull Request của bạn đã được Merge, bạn có thể chuyển về nhánh `dev` và xoá nhánh cá nhân cũ trên máy tính cho đỡ rác:
    ```bash
    git checkout dev
    git pull origin dev
    git branch -d feature/TSK-010
    ```
*   **Gõ sai tên commit?** Đừng lo, lệnh này giúp bạn sửa lời nhắn commit cuối cùng (trước khi push):
    ```bash
    git commit --amend -m "lời-nhắn-mới-chính-xác-hơn"
    ```
---
