package tests

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"book-management-api/config"
	"book-management-api/models"
	"book-management-api/routes"
	"book-management-api/utils"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"gorm.io/gorm"
)

// setupTestEnvironment menginisialisasi router dan database SQLite in-memory untuk testing
func setupTestEnvironment(t *testing.T) (*gin.Engine, *gorm.DB) {
	gin.SetMode(gin.TestMode)

	// Inisialisasi SQLite in-memory database dengan identifier unik per test
	dbPath := fmt.Sprintf("file:memdb_%s?mode=memory&cache=shared", t.Name())
	db, err := config.InitDB(dbPath)
	assert.NoError(t, err)

	r := routes.SetupRouter(db)
	return r, db
}

func TestCreateBook_Success(t *testing.T) {
	r, _ := setupTestEnvironment(t)

	payload := models.CreateBookInput{
		Title:         "Pemrograman Go Praktis",
		Author:        "Lyson Manik",
		Description:   "Buku panduan lengkap REST API dengan Gin dan GORM",
		PublishedYear: 2026,
	}
	body, _ := json.Marshal(payload)

	req, _ := http.NewRequest(http.MethodPost, "/api/books", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusCreated, w.Code)

	var res utils.APIResponse
	err := json.Unmarshal(w.Body.Bytes(), &res)
	assert.NoError(t, err)
	assert.True(t, res.Success)
	assert.Equal(t, "Buku berhasil ditambahkan", res.Message)

	dataBytes, _ := json.Marshal(res.Data)
	var createdBook models.Book
	json.Unmarshal(dataBytes, &createdBook)
	assert.Equal(t, "Pemrograman Go Praktis", createdBook.Title)
	assert.Equal(t, "Lyson Manik", createdBook.Author)
	assert.Equal(t, 2026, createdBook.PublishedYear)
}

func TestCreateBook_ValidationFailure_MissingTitle(t *testing.T) {
	r, _ := setupTestEnvironment(t)

	payload := map[string]interface{}{
		"author":         "Lyson Manik",
		"description":    "Tanpa judul buku",
		"published_year": 2026,
	}
	body, _ := json.Marshal(payload)

	req, _ := http.NewRequest(http.MethodPost, "/api/books", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)

	var res utils.APIResponse
	err := json.Unmarshal(w.Body.Bytes(), &res)
	assert.NoError(t, err)
	assert.False(t, res.Success)
	assert.Contains(t, res.Message, "wajib diisi")
}

func TestCreateBook_ValidationFailure_MissingAuthor(t *testing.T) {
	r, _ := setupTestEnvironment(t)

	payload := map[string]interface{}{
		"title":          "Belajar Pemrograman Web",
		"description":    "Tanpa penulis",
		"published_year": 2026,
	}
	body, _ := json.Marshal(payload)

	req, _ := http.NewRequest(http.MethodPost, "/api/books", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)

	var res utils.APIResponse
	err := json.Unmarshal(w.Body.Bytes(), &res)
	assert.NoError(t, err)
	assert.False(t, res.Success)
	assert.Contains(t, res.Message, "wajib diisi")
}

func TestCreateBook_ValidationFailure_EmptyWhitespace(t *testing.T) {
	r, _ := setupTestEnvironment(t)

	payload := models.CreateBookInput{
		Title:         "   ",
		Author:        "   ",
		Description:   "Hanya whitespace",
		PublishedYear: 2026,
	}
	body, _ := json.Marshal(payload)

	req, _ := http.NewRequest(http.MethodPost, "/api/books", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)

	var res utils.APIResponse
	err := json.Unmarshal(w.Body.Bytes(), &res)
	assert.NoError(t, err)
	assert.False(t, res.Success)
	assert.Contains(t, res.Message, "tidak boleh kosong")
}

func TestGetBooks_AllAndFilter(t *testing.T) {
	r, db := setupTestEnvironment(t)

	// Seed data
	b1 := models.Book{Title: "Clean Code", Author: "Robert C. Martin", PublishedYear: 2008}
	b2 := models.Book{Title: "The Pragmatic Programmer", Author: "Andrew Hunt", PublishedYear: 1999}
	b3 := models.Book{Title: "Clean Architecture", Author: "Robert C. Martin", PublishedYear: 2017}
	db.Create(&b1)
	db.Create(&b2)
	db.Create(&b3)

	// Test 1: GET All Books
	{
		req, _ := http.NewRequest(http.MethodGet, "/api/books", nil)
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
		var res utils.APIResponse
		json.Unmarshal(w.Body.Bytes(), &res)
		assert.True(t, res.Success)

		dataBytes, _ := json.Marshal(res.Data)
		var books []models.Book
		json.Unmarshal(dataBytes, &books)
		assert.Equal(t, 3, len(books))
	}

	// Test 2: GET Filter by Author (case-insensitive partial match: "martin")
	{
		req, _ := http.NewRequest(http.MethodGet, "/api/books?author=martin", nil)
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
		var res utils.APIResponse
		json.Unmarshal(w.Body.Bytes(), &res)
		assert.True(t, res.Success)

		dataBytes, _ := json.Marshal(res.Data)
		var books []models.Book
		json.Unmarshal(dataBytes, &books)
		assert.Equal(t, 2, len(books))
		assert.Equal(t, "Clean Code", books[0].Title)
		assert.Equal(t, "Clean Architecture", books[1].Title)
	}

	// Test 3: GET Filter by Author not found
	{
		req, _ := http.NewRequest(http.MethodGet, "/api/books?author=NonExistentAuthor", nil)
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
		var res utils.APIResponse
		json.Unmarshal(w.Body.Bytes(), &res)
		assert.True(t, res.Success)

		dataBytes, _ := json.Marshal(res.Data)
		var books []models.Book
		json.Unmarshal(dataBytes, &books)
		assert.Equal(t, 0, len(books))
	}
}

func TestGetBookByID(t *testing.T) {
	r, db := setupTestEnvironment(t)

	book := models.Book{Title: "Refactoring", Author: "Martin Fowler", PublishedYear: 1999}
	db.Create(&book)

	// Kasus 1: Ditemukan (200 OK)
	{
		req, _ := http.NewRequest(http.MethodGet, fmt.Sprintf("/api/books/%d", book.ID), nil)
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
		var res utils.APIResponse
		json.Unmarshal(w.Body.Bytes(), &res)
		assert.True(t, res.Success)

		dataBytes, _ := json.Marshal(res.Data)
		var foundBook models.Book
		json.Unmarshal(dataBytes, &foundBook)
		assert.Equal(t, "Refactoring", foundBook.Title)
	}

	// Kasus 2: ID tidak ditemukan (404 Not Found)
	{
		req, _ := http.NewRequest(http.MethodGet, "/api/books/9999", nil)
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusNotFound, w.Code)
		var res utils.APIResponse
		json.Unmarshal(w.Body.Bytes(), &res)
		assert.False(t, res.Success)
		assert.Contains(t, res.Message, "tidak ditemukan")
	}

	// Kasus 3: Format ID tidak valid (400 Bad Request)
	{
		req, _ := http.NewRequest(http.MethodGet, "/api/books/abc", nil)
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusBadRequest, w.Code)
		var res utils.APIResponse
		json.Unmarshal(w.Body.Bytes(), &res)
		assert.False(t, res.Success)
		assert.Contains(t, res.Message, "harus berupa angka")
	}
}

func TestUpdateBook(t *testing.T) {
	r, db := setupTestEnvironment(t)

	book := models.Book{Title: "Original Title", Author: "Original Author", PublishedYear: 2020}
	db.Create(&book)

	// Kasus 1: Update sukses (200 OK)
	{
		payload := models.UpdateBookInput{
			Title:         "Updated Title",
			Author:        "Updated Author",
			Description:   "Updated description",
			PublishedYear: 2025,
		}
		body, _ := json.Marshal(payload)

		req, _ := http.NewRequest(http.MethodPut, fmt.Sprintf("/api/books/%d", book.ID), bytes.NewBuffer(body))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
		var res utils.APIResponse
		json.Unmarshal(w.Body.Bytes(), &res)
		assert.True(t, res.Success)

		dataBytes, _ := json.Marshal(res.Data)
		var updated models.Book
		json.Unmarshal(dataBytes, &updated)
		assert.Equal(t, "Updated Title", updated.Title)
		assert.Equal(t, "Updated Author", updated.Author)
		assert.Equal(t, 2025, updated.PublishedYear)
	}

	// Kasus 2: Validasi gagal - title kosong (400 Bad Request)
	{
		payload := map[string]interface{}{
			"title":  "",
			"author": "Author Only",
		}
		body, _ := json.Marshal(payload)

		req, _ := http.NewRequest(http.MethodPut, fmt.Sprintf("/api/books/%d", book.ID), bytes.NewBuffer(body))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusBadRequest, w.Code)
		var res utils.APIResponse
		json.Unmarshal(w.Body.Bytes(), &res)
		assert.False(t, res.Success)
		assert.Contains(t, res.Message, "wajib diisi")
	}

	// Kasus 3: ID tidak ditemukan (404 Not Found)
	{
		payload := models.UpdateBookInput{
			Title:  "Valid Title",
			Author: "Valid Author",
		}
		body, _ := json.Marshal(payload)

		req, _ := http.NewRequest(http.MethodPut, "/api/books/8888", bytes.NewBuffer(body))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusNotFound, w.Code)
		var res utils.APIResponse
		json.Unmarshal(w.Body.Bytes(), &res)
		assert.False(t, res.Success)
		assert.Contains(t, res.Message, "tidak ditemukan")
	}
}

func TestDeleteBook(t *testing.T) {
	r, db := setupTestEnvironment(t)

	book := models.Book{Title: "Book To Delete", Author: "Some Author", PublishedYear: 2021}
	db.Create(&book)

	// Kasus 1: Hapus sukses (200 OK)
	{
		req, _ := http.NewRequest(http.MethodDelete, fmt.Sprintf("/api/books/%d", book.ID), nil)
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
		var res utils.APIResponse
		json.Unmarshal(w.Body.Bytes(), &res)
		assert.True(t, res.Success)
		assert.Contains(t, res.Message, "berhasil dihapus")

		// Verifikasi dari DB bahwa data sudah terhapus
		var checkBook models.Book
		err := db.First(&checkBook, book.ID).Error
		assert.Equal(t, gorm.ErrRecordNotFound, err)
	}

	// Kasus 2: Hapus ID yang tidak ada (404 Not Found)
	{
		req, _ := http.NewRequest(http.MethodDelete, "/api/books/9999", nil)
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusNotFound, w.Code)
		var res utils.APIResponse
		json.Unmarshal(w.Body.Bytes(), &res)
		assert.False(t, res.Success)
		assert.Contains(t, res.Message, "tidak ditemukan")
	}
}

func TestNoRoute_NotFound(t *testing.T) {
	r, _ := setupTestEnvironment(t)

	req, _ := http.NewRequest(http.MethodGet, "/api/unknown-endpoint", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusNotFound, w.Code)
	var res utils.APIResponse
	json.Unmarshal(w.Body.Bytes(), &res)
	assert.False(t, res.Success)
	assert.Contains(t, res.Message, "404 Not Found")
}
