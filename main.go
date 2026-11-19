package main

import (
	"fmt"
	"log"
	"os"

	"book-management-api/config"
	"book-management-api/routes"
)

func main() {
	// Inisialisasi Database Relasional (SQLite)
	dbPath := os.Getenv("DB_PATH")
	if dbPath == "" {
		dbPath = "books.db"
	}

	db, err := config.InitDB(dbPath)
	if err != nil {
		log.Fatalf("Gagal inisialisasi basis data: %v", err)
	}

	// Inisialisasi Router Gin
	r := routes.SetupRouter(db)

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	addr := fmt.Sprintf(":%s", port)
	log.Printf("Server Book Management API berjalan pada http://localhost%s", addr)
	if err := r.Run(addr); err != nil {
		log.Fatalf("Gagal menjalankan server: %v", err)
	}
}
