package models

import "github.com/jinzhu/gorm"

// StoreURL ...
type StoreURL struct {
	gorm.Model

	ID           int    `json:"id"`
	GeneratedURL string `json:"generatedURL"`
	OriginalURL  string `json:"originalURL"`
	CreatedAt    string `json:"createdAt"`
	UpdatedAt    string `json:"updatedAt"`
}
