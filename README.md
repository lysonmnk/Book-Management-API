# Book Management API

RESTful API untuk manajemen data buku yang dibangun menggunakan bahasa pemrograman **Go (Golang)**, framework web **Gin**, dan ORM **GORM** dengan database relasional **SQLite**.

---

## 📋 Daftar Fitur

1. **Operasi CRUD Lengkap**:
   - `POST /api/books`: Menambahkan buku baru.
   - `GET /api/books`: Mengambil seluruh daftar buku.
   - `GET /api/books/:id`: Mengambil detail satu buku berdasarkan ID.
   - `PUT /api/books/:id`: Memperbarui data buku berdasarkan ID.
   - `DELETE /api/books/:id`: Menghapus data buku berdasarkan ID.
2. **Validasi Input**:
   - Field `title` dan `author` wajib diisi (*required*).
   - Validasi terhadap string kosong atau hanya karakter spasi (*whitespace*). Jika tidak valid, mengembalikan status **400 Bad Request**.
3. **Pencarian / Filter Penulis**:
   - Endpoint `GET /api/books?author=...` untuk memfilter buku berdasarkan nama penulis secara fleksibel (*case-insensitive partial match*).
4. **Format Response Konsisten**:
   - Format response JSON yang terstandar baik untuk response berhasil (**200**, **201**) maupun response error (**400**, **404**, **500**).
   - Dilengkapi handler `NoRoute` terstandar untuk URL yang tidak terdaftar.
5. **Bonus Database Relasional**:
   - Menggunakan database relasional **SQLite** (`books.db`) via GORM (`github.com/glebarez/sqlite` pure-Go) sehingga data tersimpan secara persisten dan mendukung transaksi SQL relasional tanpa memerlukan instalasi database server eksternal.
6. **Automated Unit & Integration Testing**:
   - Pengujian terotomatisasi mencakup seluruh alur sukses, validasi, dan penanganan error.

---

## 📂 Struktur Proyek

```
Book Management API/
├── config/
│   └── database.go          # Inisialisasi GORM & auto-migrasi tabel SQLite
├── controllers/
│   └── book_controller.go   # Handler logika bisnis CRUD, filter, dan validasi
├── models/
│   └── book.go              # Model struct GORM Book & DTO Create/Update
├── routes/
│   └── routes.go            # Konfigurasi router Gin, middleware CORS, dan route handler
├── tests/
│   └── api_test.go          # Suite automated unit/integration tests
├── utils/
│   └── response.go          # Helper response JSON konsisten
├── .gitignore               # Konfigurasi file yang diabaikan oleh Git
├── go.mod                   # File definisi modul Go
├── go.sum                   # Checksum dependensi Go
├── main.go                  # Entry point aplikasi
├── README.md                # Dokumentasi utama proyek & panduan penggunaan
└── TESTING.md               # Bukti dan dokumentasi pengujian lengkap endpoint
```

---

## 🛠️ Prasyarat (Prerequisites)

- **Go (Golang)**: Versi 1.20 atau yang lebih baru (dapat diperiksa dengan `go version`).
- **Git**

---

## 🚀 Panduan Menjalankan Aplikasi

### 1. Clone Repository (Jika belum)
```bash
git clone https://github.com/lysonmnk/Book-Management-API.git
cd "Book Management API"
```

### 2. Unduh Dependensi
```bash
go mod tidy
```

### 3. Jalankan Server
Secara default server akan berjalan pada port `8080` dan membuat file database `books.db` secara otomatis:
```bash
go run main.go
```

Output terminal saat server berjalan:
```text
[GIN-debug] [WARNING] Creating an Engine instance with the Logger and Recovery middleware already attached.
...
Server Book Management API berjalan pada http://localhost:8080
```

*(Opsional)* Anda dapat mengubah port atau path database menggunakan *environment variables*:
```powershell
$env:PORT="3000"
$env:DB_PATH="custom_books.db"
go run main.go
```

---

## 🧪 Menjalankan Automated Tests

Aplikasi dilengkapi dengan test suite lengkap menggunakan in-memory SQLite database sehingga pengujian berjalan cepat, independen, dan terisolasi.

Jalankan perintah berikut:
```bash
go test -v ./...
```

Contoh output hasil pengujian yang sukses:
```text
=== RUN   TestCreateBook_Success
--- PASS: TestCreateBook_Success (0.00s)
=== RUN   TestCreateBook_ValidationFailure_MissingTitle
--- PASS: TestCreateBook_ValidationFailure_MissingTitle (0.00s)
=== RUN   TestCreateBook_ValidationFailure_MissingAuthor
--- PASS: TestCreateBook_ValidationFailure_MissingAuthor (0.00s)
=== RUN   TestCreateBook_ValidationFailure_EmptyWhitespace
--- PASS: TestCreateBook_ValidationFailure_EmptyWhitespace (0.00s)
=== RUN   TestGetBooks_AllAndFilter
--- PASS: TestGetBooks_AllAndFilter (0.00s)
=== RUN   TestGetBookByID
--- PASS: TestGetBookByID (0.00s)
=== RUN   TestUpdateBook
--- PASS: TestUpdateBook (0.00s)
=== RUN   TestDeleteBook
--- PASS: TestDeleteBook (0.00s)
=== RUN   TestNoRoute_NotFound
--- PASS: TestNoRoute_NotFound (0.00s)
PASS
ok      book-management-api/tests       0.397s
```

---

## 📖 Dokumentasi Endpoint API

Base URL: `http://localhost:8080`

| Method | Endpoint | Deskripsi | Status Code Sukses | Status Code Error |
| :--- | :--- | :--- | :---: | :---: |
| `GET` | `/` | Cek status server API | `200 OK` | - |
| `POST` | `/api/books` | Menambahkan buku baru | `201 Created` | `400 Bad Request` |
| `GET` | `/api/books` | Mengambil semua buku | `200 OK` | `500 Internal Error` |
| `GET` | `/api/books?author=...` | Filter buku berdasarkan nama penulis | `200 OK` | `500 Internal Error` |
| `GET` | `/api/books/:id` | Mengambil detail satu buku berdasarkan ID | `200 OK` | `400 Bad Request`, `404 Not Found` |
| `PUT` | `/api/books/:id` | Memperbarui data buku | `200 OK` | `400 Bad Request`, `404 Not Found` |
| `DELETE`| `/api/books/:id` | Menghapus buku | `200 OK` | `400 Bad Request`, `404 Not Found` |

---

## 📦 Standarisasi Format Response

### 1. Format Response Sukses (200 / 201)
```json
{
  "success": true,
  "message": "Pesan keberhasilan operasi",
  "data": { ... }
}
```

### 2. Format Response Error (400 / 404 / 500)
```json
{
  "success": false,
  "message": "Deskripsi kesalahan yang terjadi"
}
```

---

## 💡 Contoh Penggunaan dengan cURL

### 1. Tambah Buku (POST)
```bash
curl -X POST http://localhost:8080/api/books \
  -H "Content-Type: application/json" \
  -d '{
    "title": "Clean Architecture",
    "author": "Robert C. Martin",
    "description": "A Craftsman'\''s Guide to Software Structure and Design",
    "published_year": 2017
  }'
```

### 2. Ambil Semua Buku (GET)
```bash
curl -X GET http://localhost:8080/api/books
```

### 3. Filter Buku Berdasarkan Penulis (GET Query)
```bash
curl -X GET "http://localhost:8080/api/books?author=Robert"
```

### 4. Ambil Detail Buku Berdasarkan ID (GET)
```bash
curl -X GET http://localhost:8080/api/books/1
```

### 5. Update Buku (PUT)
```bash
curl -X PUT http://localhost:8080/api/books/1 \
  -H "Content-Type: application/json" \
  -d '{
    "title": "Clean Architecture (2nd Edition)",
    "author": "Robert C. Martin",
    "description": "Edisi revisi terbaru",
    "published_year": 2020
  }'
```

### 6. Hapus Buku (DELETE)
```bash
curl -X DELETE http://localhost:8080/api/books/1
```

> Untuk bukti pengujian mendalam beserta tangkapan format respon dari masing-masing kasus uji, silakan lihat [TESTING.md](TESTING.md).

---

## 👤 Pengembang

- **Repositori**: [lysonmnk/Book-Management-API](https://github.com/lysonmnk/Book-Management-API)
