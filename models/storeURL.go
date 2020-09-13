package models

// StoreURL ...
type StoreURL struct {
	ID           int    `json:"id"`
	GeneratedURL string `json:"generatedURL"`
	OriginalURL  string `json:"originalURL"`
	CreatedAt    string `json:"createdAt"`
	UpdatedAt    string `json:"updatedAt"`
}
