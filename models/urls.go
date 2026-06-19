package models

import "time"

type URL struct {
	ID           int       `json:"id"`
	OriginalURL  string    `json:"originalURL"`
	ShortenURL   string    `json:"shortURL"`
	CreationDate time.Time `json:"creationDate"`
}
