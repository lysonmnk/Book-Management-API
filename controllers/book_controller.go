package controllers

import (
	"net/http"
	"strconv"
	"strings"

	"book-management-api/models"
	"book-management-api/utils"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// BookController menangani request HTTP untuk resource buku
type BookController struct {
	DB *gorm.DB
}

// NewBookController mengembalikan instance baru dari BookController
func NewBookController(db *gorm.DB) *BookController {
	return &BookController{DB: db}
}

// CreateBook membuat buku baru (POST /api/books)
func (ctrl *BookController) CreateBook(c *gin.Context) {
	var input models.CreateBookInput

	if err := c.ShouldBindJSON(&input); err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "Field 'title' dan 'author' wajib diisi")
		return
	}

	title := strings.TrimSpace(input.Title)
	author := strings.TrimSpace(input.Author)

	if title == "" || author == "" {
		utils.ErrorResponse(c, http.StatusBadRequest, "Field 'title' dan 'author' wajib diisi dan tidak boleh kosong")
		return
	}

	book := models.Book{
		Title:         title,
		Author:        author,
		Description:   strings.TrimSpace(input.Description),
		PublishedYear: input.PublishedYear,
	}

	if err := ctrl.DB.Create(&book).Error; err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, "Gagal menyimpan buku ke database: "+err.Error())
		return
	}

	utils.SuccessResponse(c, http.StatusCreated, "Buku berhasil ditambahkan", book)
}

// GetBooks mengambil semua buku dengan opsi filter berdasarkan author (GET /api/books?author=...)
func (ctrl *BookController) GetBooks(c *gin.Context) {
	var books = make([]models.Book, 0)
	authorFilter := strings.TrimSpace(c.Query("author"))

	query := ctrl.DB.Model(&models.Book{})
	if authorFilter != "" {
		// Filter case-insensitive partial match
		query = query.Where("LOWER(author) LIKE LOWER(?)", "%"+authorFilter+"%")
	}

	if err := query.Order("id asc").Find(&books).Error; err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, "Gagal mengambil data buku: "+err.Error())
		return
	}

	utils.SuccessResponse(c, http.StatusOK, "Data buku berhasil diambil", books)
}

// GetBookByID mengambil detail satu buku berdasarkan ID (GET /api/books/:id)
func (ctrl *BookController) GetBookByID(c *gin.Context) {
	idParam := c.Param("id")
	id, err := strconv.ParseUint(idParam, 10, 64)
	if err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "Format ID buku tidak valid, ID harus berupa angka")
		return
	}

	var book models.Book
	if err := ctrl.DB.First(&book, id).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			utils.ErrorResponse(c, http.StatusNotFound, "Buku dengan ID "+idParam+" tidak ditemukan")
			return
		}
		utils.ErrorResponse(c, http.StatusInternalServerError, "Gagal mengambil data buku: "+err.Error())
		return
	}

	utils.SuccessResponse(c, http.StatusOK, "Detail buku berhasil diambil", book)
}

// UpdateBook memperbarui data buku berdasarkan ID (PUT /api/books/:id)
func (ctrl *BookController) UpdateBook(c *gin.Context) {
	idParam := c.Param("id")
	id, err := strconv.ParseUint(idParam, 10, 64)
	if err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "Format ID buku tidak valid, ID harus berupa angka")
		return
	}

	var book models.Book
	if err := ctrl.DB.First(&book, id).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			utils.ErrorResponse(c, http.StatusNotFound, "Buku dengan ID "+idParam+" tidak ditemukan")
			return
		}
		utils.ErrorResponse(c, http.StatusInternalServerError, "Gagal mengambil data buku: "+err.Error())
		return
	}

	var input models.UpdateBookInput
	if err := c.ShouldBindJSON(&input); err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "Field 'title' dan 'author' wajib diisi")
		return
	}

	title := strings.TrimSpace(input.Title)
	author := strings.TrimSpace(input.Author)

	if title == "" || author == "" {
		utils.ErrorResponse(c, http.StatusBadRequest, "Field 'title' dan 'author' wajib diisi dan tidak boleh kosong")
		return
	}

	book.Title = title
	book.Author = author
	book.Description = strings.TrimSpace(input.Description)
	book.PublishedYear = input.PublishedYear

	if err := ctrl.DB.Save(&book).Error; err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, "Gagal memperbarui data buku: "+err.Error())
		return
	}

	utils.SuccessResponse(c, http.StatusOK, "Data buku berhasil diperbarui", book)
}

// DeleteBook menghapus data buku berdasarkan ID (DELETE /api/books/:id)
func (ctrl *BookController) DeleteBook(c *gin.Context) {
	idParam := c.Param("id")
	id, err := strconv.ParseUint(idParam, 10, 64)
	if err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "Format ID buku tidak valid, ID harus berupa angka")
		return
	}

	var book models.Book
	if err := ctrl.DB.First(&book, id).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			utils.ErrorResponse(c, http.StatusNotFound, "Buku dengan ID "+idParam+" tidak ditemukan")
			return
		}
		utils.ErrorResponse(c, http.StatusInternalServerError, "Gagal mengambil data buku: "+err.Error())
		return
	}

	if err := ctrl.DB.Delete(&book).Error; err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, "Gagal menghapus buku: "+err.Error())
		return
	}

	utils.SuccessResponse(c, http.StatusOK, "Buku dengan ID "+idParam+" berhasil dihapus", nil)
}
