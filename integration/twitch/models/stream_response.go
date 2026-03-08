package models

import "time"

type StreamData struct {
	ID           string    `json:"id"`
	UserID       string    `json:"user_id"`
	UserLogin    string    `json:"user_login"`
	UserName     string    `json:"user_name"`
	GameID       string    `json:"game_id"`
	GameName     string    `json:"game_name"`
	Type         string    `json:"type"`
	Title        string    `json:"title"`
	ViewerCount  int       `json:"viewer_count"`
	StartedAt    time.Time `json:"started_at"`
	Language     string    `json:"language"`
	ThumbnailURL string    `json:"thumbnail_url"`
	TagIds       []any     `json:"tag_ids"`
	Tags         []string  `json:"tags"`
	IsMature     bool      `json:"is_mature"`
}

type StreamResponse struct {
	Data []StreamData `json:"data"`
}
