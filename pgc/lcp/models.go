package lcp

import "time"

type CacheData interface {
	[]GitHubRepository | []HackerNews
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

type HackerNews struct {
	ID    int    `json:"id"`
	Title string `json:"title"`
	By    string `json:"by"`
	URL   string `json:"url"`
	Score int    `json:"score"`
	Time  int64  `json:"time"`
}
