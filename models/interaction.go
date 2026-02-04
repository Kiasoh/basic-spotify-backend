package models

import "time"

type Interaction struct {
	UserID    int       `json:"user_id"`
	SongID    int       `json:"song_id"`
	Type      string    `json:"type"`
	CreatedAt time.Time `json:"created_at"`
}
