# Bukti Pengujian Endpoint (Testing Evidence)
**Book Management API**  
Framework: **Gin (Golang)** | Database: **Relasional SQLite via GORM**

Dokumen ini berisi dokumentasi dan bukti pengujian lengkap seluruh endpoint API sesuai spesifikasi kebutuhan.

---

## 1. Bukti Eksekusi Automated Test Suite (`go test`)

Seluruh skenario pengujian diimplementasikan menggunakan `testing`, `net/http/httptest`, dan `testify/assert` dengan in-memory SQLite database (`tests/api_test.go`).

Jalankan perintah:
```powershell
go test -v ./...
```

### Output Eksekusi Pengujian:
```text
=== RUN   TestCreateBook_Success
[GIN] 2026/11/19 - 14:39:58 | 201 |   1.24ms |                 | POST     "/api/books"
--- PASS: TestCreateBook_Success (0.00s)
=== RUN   TestCreateBook_ValidationFailure_MissingTitle
[GIN] 2026/11/19 - 14:39:58 | 400 |       0s |                 | POST     "/api/books"
--- PASS: TestCreateBook_ValidationFailure_MissingTitle (0.00s)
=== RUN   TestCreateBook_ValidationFailure_MissingAuthor
[GIN] 2026/11/19 - 14:39:58 | 400 |       0s |                 | POST     "/api/books"
--- PASS: TestCreateBook_ValidationFailure_MissingAuthor (0.00s)
=== RUN   TestCreateBook_ValidationFailure_EmptyWhitespace
[GIN] 2026/11/19 - 14:39:58 | 400 |       0s |                 | POST     "/api/books"
--- PASS: TestCreateBook_ValidationFailure_EmptyWhitespace (0.00s)
=== RUN   TestGetBooks_AllAndFilter
[GIN] 2026/11/19 - 14:39:58 | 200 |       0s |                 | GET      "/api/books"
[GIN] 2026/11/19 - 14:39:58 | 200 |       0s |                 | GET      "/api/books?author=martin"
[GIN] 2026/11/19 - 14:39:58 | 200 |       0s |                 | GET      "/api/books?author=NonExistentAuthor"
--- PASS: TestGetBooks_AllAndFilter (0.00s)
=== RUN   TestGetBookByID
[GIN] 2026/11/19 - 14:39:58 | 200 |  519.9µs |                 | GET      "/api/books/1"
[GIN] 2026/11/19 - 14:39:58 | 404 |       0s |                 | GET      "/api/books/9999"
[GIN] 2026/11/19 - 14:39:58 | 400 |       0s |                 | GET      "/api/books/abc"
--- PASS: TestGetBookByID (0.00s)
=== RUN   TestUpdateBook
[GIN] 2026/11/19 - 14:39:58 | 200 |       0s |                 | PUT      "/api/books/1"
[GIN] 2026/11/19 - 14:39:58 | 400 |       0s |                 | PUT      "/api/books/1"
[GIN] 2026/11/19 - 14:39:58 | 404 |  517.9µs |                 | PUT      "/api/books/8888"
--- PASS: TestUpdateBook (0.00s)
=== RUN   TestDeleteBook
[GIN] 2026/11/19 - 14:39:58 | 200 |  513.3µs |                 | DELETE   "/api/books/1"
[GIN] 2026/11/19 - 14:39:58 | 404 |       0s |                 | DELETE   "/api/books/9999"
--- PASS: TestDeleteBook (0.00s)
=== RUN   TestNoRoute_NotFound
[GIN] 2026/11/19 - 14:39:58 | 404 |       0s |                 | GET      "/api/unknown-endpoint"
--- PASS: TestNoRoute_NotFound (0.00s)
PASS
ok  	book-management-api/tests	0.397s
```

---

## 2. Rincian Pengujian Endpoint HTTP

### A. Operasi Create (`POST /api/books`)

#### 1. Berhasil Menambahkan Buku
- **Method & URL**: `POST http://localhost:8080/api/books`
- **Request Headers**: `Content-Type: application/json`
- **Request Body**:
```json
{
  "title": "Clean Architecture",
  "author": "Robert C. Martin",
  "description": "Panduan praktis arsitektur perangkat lunak yang kokoh",
  "published_year": 2017
}
```
- **cURL**:
```bash
curl -X POST http://localhost:8080/api/books \
  -H "Content-Type: application/json" \
  -d '{"title":"Clean Architecture","author":"Robert C. Martin","description":"Panduan praktis arsitektur perangkat lunak yang kokoh","published_year":2017}'
```
- **HTTP Status**: `201 Created`
- **Response Body**:
```json
{
  "success": true,
  "message": "Buku berhasil ditambahkan",
  "data": {
    "id": 1,
    "title": "Clean Architecture",
    "author": "Robert C. Martin",
    "description": "Panduan praktis arsitektur perangkat lunak yang kokoh",
    "published_year": 2017,
    "created_at": "2026-11-19T14:40:00Z",
    "updated_at": "2026-11-19T14:40:00Z"
  }
}
```

---

#### 2. Gagal Validasi: Field `title` Kosong
- **Method & URL**: `POST http://localhost:8080/api/books`
- **Request Body**:
```json
{
  "author": "Robert C. Martin",
  "published_year": 2017
}
```
- **cURL**:
```bash
curl -X POST http://localhost:8080/api/books \
  -H "Content-Type: application/json" \
  -d '{"author":"Robert C. Martin","published_year":2017}'
```
- **HTTP Status**: `400 Bad Request`
- **Response Body**:
```json
{
  "success": false,
  "message": "Field 'title' dan 'author' wajib diisi"
}
```

---

#### 3. Gagal Validasi: Field `author` Kosong / Whitespace
- **Method & URL**: `POST http://localhost:8080/api/books`
- **Request Body**:
```json
{
  "title": "Clean Code",
  "author": "   "
}
```
- **cURL**:
```bash
curl -X POST http://localhost:8080/api/books \
  -H "Content-Type: application/json" \
  -d '{"title":"Clean Code","author":"   "}'
```
- **HTTP Status**: `400 Bad Request`
- **Response Body**:
```json
{
  "success": false,
  "message": "Field 'title' dan 'author' wajib diisi dan tidak boleh kosong"
}
```

---

### B. Operasi Read All & Filter Author (`GET /api/books`)

#### 1. Mengambil Semua Buku
- **Method & URL**: `GET http://localhost:8080/api/books`
- **cURL**:
```bash
curl -X GET http://localhost:8080/api/books
```
- **HTTP Status**: `200 OK`
- **Response Body**:
```json
{
  "success": true,
  "message": "Data buku berhasil diambil",
  "data": [
    {
      "id": 1,
      "title": "Clean Architecture",
      "author": "Robert C. Martin",
      "description": "Panduan praktis arsitektur perangkat lunak yang kokoh",
      "published_year": 2017,
      "created_at": "2026-11-19T14:40:00Z",
      "updated_at": "2026-11-19T14:40:00Z"
    },
    {
      "id": 2,
      "title": "The Pragmatic Programmer",
      "author": "Andrew Hunt",
      "description": "Your Journey to Mastery",
      "published_year": 1999,
      "created_at": "2026-11-19T14:40:01Z",
      "updated_at": "2026-11-19T14:40:01Z"
    }
  ]
}
```

---

#### 2. Filter Berdasarkan Penulis (`GET /api/books?author=...`)
- **Method & URL**: `GET http://localhost:8080/api/books?author=martin`
- **cURL**:
```bash
curl -X GET "http://localhost:8080/api/books?author=martin"
```
- **HTTP Status**: `200 OK`
- **Response Body**:
```json
{
  "success": true,
  "message": "Data buku berhasil diambil",
  "data": [
    {
      "id": 1,
      "title": "Clean Architecture",
      "author": "Robert C. Martin",
      "description": "Panduan praktis arsitektur perangkat lunak yang kokoh",
      "published_year": 2017,
      "created_at": "2026-11-19T14:40:00Z",
      "updated_at": "2026-11-19T14:40:00Z"
    }
  ]
}
```

---

### C. Operasi Read One (`GET /api/books/:id`)

#### 1. Buku Ditemukan
- **Method & URL**: `GET http://localhost:8080/api/books/1`
- **cURL**:
```bash
curl -X GET http://localhost:8080/api/books/1
```
- **HTTP Status**: `200 OK`
- **Response Body**:
```json
{
  "success": true,
  "message": "Detail buku berhasil diambil",
  "data": {
    "id": 1,
    "title": "Clean Architecture",
    "author": "Robert C. Martin",
    "description": "Panduan praktis arsitektur perangkat lunak yang kokoh",
    "published_year": 2017,
    "created_at": "2026-11-19T14:40:00Z",
    "updated_at": "2026-11-19T14:40:00Z"
  }
}
```

---

#### 2. Buku Tidak Ditemukan (404 Not Found)
- **Method & URL**: `GET http://localhost:8080/api/books/999`
- **cURL**:
```bash
curl -X GET http://localhost:8080/api/books/999
```
- **HTTP Status**: `404 Not Found`
- **Response Body**:
```json
{
  "success": false,
  "message": "Buku dengan ID 999 tidak ditemukan"
}
```

---

#### 3. Format ID Tidak Valid (400 Bad Request)
- **Method & URL**: `GET http://localhost:8080/api/books/abc`
- **cURL**:
```bash
curl -X GET http://localhost:8080/api/books/abc
```
- **HTTP Status**: `400 Bad Request`
- **Response Body**:
```json
{
  "success": false,
  "message": "Format ID buku tidak valid, ID harus berupa angka"
}
```

---

### D. Operasi Update (`PUT /api/books/:id`)

#### 1. Berhasil Memperbarui Buku
- **Method & URL**: `PUT http://localhost:8080/api/books/1`
- **Request Body**:
```json
{
  "title": "Clean Architecture: A Craftsman's Guide",
  "author": "Robert C. Martin",
  "description": "Edisi revisi terbaru",
  "published_year": 2018
}
```
- **cURL**:
```bash
curl -X PUT http://localhost:8080/api/books/1 \
  -H "Content-Type: application/json" \
  -d '{"title":"Clean Architecture: A Craftsman'\''s Guide","author":"Robert C. Martin","description":"Edisi revisi terbaru","published_year":2018}'
```
- **HTTP Status**: `200 OK`
- **Response Body**:
```json
{
  "success": true,
  "message": "Data buku berhasil diperbarui",
  "data": {
    "id": 1,
    "title": "Clean Architecture: A Craftsman's Guide",
    "author": "Robert C. Martin",
    "description": "Edisi revisi terbaru",
    "published_year": 2018,
    "created_at": "2026-11-19T14:40:00Z",
    "updated_at": "2026-11-19T14:45:00Z"
  }
}
```

---

#### 2. Gagal Update: ID Tidak Ditemukan (404 Not Found)
- **Method & URL**: `PUT http://localhost:8080/api/books/999`
- **Request Body**:
```json
{
  "title": "Buku Baru",
  "author": "Penulis Baru"
}
```
- **HTTP Status**: `404 Not Found`
- **Response Body**:
```json
{
  "success": false,
  "message": "Buku dengan ID 999 tidak ditemukan"
}
```

---

### E. Operasi Delete (`DELETE /api/books/:id`)

#### 1. Berhasil Menghapus Buku
- **Method & URL**: `DELETE http://localhost:8080/api/books/1`
- **cURL**:
```bash
curl -X DELETE http://localhost:8080/api/books/1
```
- **HTTP Status**: `200 OK`
- **Response Body**:
```json
{
  "success": true,
  "message": "Buku dengan ID 1 berhasil dihapus"
}
```

---

#### 2. Gagal Hapus: ID Tidak Ditemukan (404 Not Found)
- **Method & URL**: `DELETE http://localhost:8080/api/books/999`
- **cURL**:
```bash
curl -X DELETE http://localhost:8080/api/books/999
```
- **HTTP Status**: `404 Not Found`
- **Response Body**:
```json
{
  "success": false,
  "message": "Buku dengan ID 999 tidak ditemukan"
}
```

---

### F. Handler NoRoute (404 Not Found)

#### Mengakses Endpoint yang Tidak Terdaftar
- **Method & URL**: `GET http://localhost:8080/api/unknown-route`
- **cURL**:
```bash
curl -X GET http://localhost:8080/api/unknown-route
```
- **HTTP Status**: `404 Not Found`
- **Response Body**:
```json
{
  "success": false,
  "message": "Endpoint tidak ditemukan (404 Not Found)"
}
```
