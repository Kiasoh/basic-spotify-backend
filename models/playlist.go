package models

import "time"

type Playlist struct {
	ID         int       `json:"id"`
	Name       string    `json:"name"`
	OwnerID    int       `json:"owner_id"`
	Modifyable bool      `json:"modifyable"`
	CreatedAt  time.Time `json:"created_at"`
}
