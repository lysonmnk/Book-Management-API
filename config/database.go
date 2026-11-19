package config

import (
	"log"

	"book-management-api/models"

	"github.com/glebarez/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

// DB adalah instance global database GORM
var DB *gorm.DB

// InitDB menginisialisasi koneksi database SQLite relasional dan melakukan auto-migrasi tabel
func InitDB(dbPath string) (*gorm.DB, error) {
	if dbPath == "" {
		dbPath = "books.db"
	}

	db, err := gorm.Open(sqlite.Open(dbPath), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Silent),
	})
	if err != nil {
		log.Printf("Gagal membuka koneksi database (%s): %v", dbPath, err)
		return nil, err
	}

	// Auto-migrasi tabel Book
	if err := db.AutoMigrate(&models.Book{}); err != nil {
		log.Printf("Gagal migrasi skema tabel Book: %v", err)
		return nil, err
	}

	DB = db
	return db, nil
}
