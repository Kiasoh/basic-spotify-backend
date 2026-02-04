package models

import "time"

type Playlist struct {
	ID          int       `json:"id"`
	Name        string    `json:"name"`
	Description *string   `json:"description,omitempty"`
	OwnerID     int       `json:"owner_id"`
	Modifyable  bool      `json:"modifyable"`
	CreatedAt   time.Time `json:"created_at"`
}
