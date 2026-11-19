package models

import "time"

// Book merepresentasikan entitas buku dalam database relasional
type Book struct {
	ID            uint      `json:"id" gorm:"primaryKey;autoIncrement"`
	Title         string    `json:"title" gorm:"not null"`
	Author        string    `json:"author" gorm:"not null;index"`
	Description   string    `json:"description" gorm:"type:text"`
	PublishedYear int       `json:"published_year"`
	CreatedAt     time.Time `json:"created_at"`
	UpdatedAt     time.Time `json:"updated_at"`
}

// CreateBookInput mendefinisikan payload untuk pembuatan buku baru
type CreateBookInput struct {
	Title         string `json:"title" binding:"required"`
	Author        string `json:"author" binding:"required"`
	Description   string `json:"description"`
	PublishedYear int    `json:"published_year"`
}

// UpdateBookInput mendefinisikan payload untuk pembaruan data buku
type UpdateBookInput struct {
	Title         string `json:"title" binding:"required"`
	Author        string `json:"author" binding:"required"`
	Description   string `json:"description"`
	PublishedYear int    `json:"published_year"`
}
