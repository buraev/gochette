package lcp

import "time"

type CacheData interface {
	[]GitHubRepository
}

type GitHubRepository struct {
	Name          string    `json:"name"`
	Owner         string    `json:"owner"`
	Language      string    `json:"language"`
	LanguageColor string    `json:"language_color"`
	Description   string    `json:"description"`
	UpdatedAt     time.Time `json:"updated_at"`
	ID            string    `json:"id"`
	URL           string    `json:"url"`
}
