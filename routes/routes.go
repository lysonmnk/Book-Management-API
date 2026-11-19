package routes

import (
	"net/http"

	"book-management-api/controllers"
	"book-management-api/utils"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// SetupRouter mengonfigurasi engine router Gin beserta route dan middleware
func SetupRouter(db *gorm.DB) *gin.Engine {
	r := gin.Default()

	// Middleware CORS sederhana agar dapat diakses dari frontend (React/Svelte/Postman dsb.)
	r.Use(func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, accept, origin, Cache-Control, X-Requested-With")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS, GET, PUT, DELETE")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(http.StatusNoContent)
			return
		}

		c.Next()
	})

	// Root endpoint
	r.GET("/", func(c *gin.Context) {
		utils.SuccessResponse(c, http.StatusOK, "Selamat datang di Book Management API", gin.H{
			"version": "1.0.0",
			"docs":    "Gunakan endpoint /api/books untuk mengelola buku",
		})
	})

	// Inisialisasi controller
	bookCtrl := controllers.NewBookController(db)

	// Routing grup /api
	api := r.Group("/api")
	{
		books := api.Group("/books")
		{
			books.POST("", bookCtrl.CreateBook)
			books.GET("", bookCtrl.GetBooks)
			books.GET("/:id", bookCtrl.GetBookByID)
			books.PUT("/:id", bookCtrl.UpdateBook)
			books.DELETE("/:id", bookCtrl.DeleteBook)
		}
	}

	// Handler konsisten untuk rute yang tidak ditemukan (404)
	r.NoRoute(func(c *gin.Context) {
		utils.ErrorResponse(c, http.StatusNotFound, "Endpoint tidak ditemukan (404 Not Found)")
	})

	return r
}
