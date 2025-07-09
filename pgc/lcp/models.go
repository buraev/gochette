package lcp

import "time"

type CacheData interface {
	[]GitHubRepository | []HackerNews | []SteamGame
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

type SteamGame struct {
	Name                string              `json:"name"`
	AppID               int32               `json:"app_id"`
	IconURL             string              `json:"icon_url"`
	RTimeLastPlayed     time.Time           `json:"rtime_last_played"`
	PlaytimeForever     int32               `json:"playtime_forever"`
	URL                 string              `json:"url"`
	HeaderURL           string              `json:"header_url"`
	LibraryHeroURL      string              `json:"library_hero_url"`
	LibraryHeroLogoURL  string              `json:"library_hero_logo_url"`
	AchievementProgress *float32            `json:"achievement_progress"`
	Achievements        *[]SteamAchievement `json:"achievements"`
}

type SteamAchievement struct {
	ApiName     string     `json:"api_name"`
	Achieved    bool       `json:"achieved"`
	Icon        string     `json:"icon"`
	DisplayName string     `json:"display_name"`
	Description *string    `json:"description"`
	UnlockTime  *time.Time `json:"unlock_time"`
}
