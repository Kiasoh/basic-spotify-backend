package models

import "time"

type TrackPlaylist struct {
	PlaylistID int       `json:"playlist_id"`
	TrackID    string    `json:"track_id"`
	CreatedAt  time.Time `json:"created_at"`
}
