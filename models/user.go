package models

import "time"

type User struct {
	ID               int       `json:"id"`
	Username         string    `json:"username"`
	Password         Password  `json:"-"`
	RecommPlaylistID int       `json:"recomm_playlist_id"`
	CreatedAt        time.Time `json:"created_at"`
}
