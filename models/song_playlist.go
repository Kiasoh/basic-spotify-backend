package models

import "time"

type SongPlaylist struct {
	PlaylistID int       `json:"playlist_id"`
	SongID     int       `json:"song_id"`
	CreatedAt  time.Time `json:"created_at"`
}
