package models

import "time"

type Song struct {
	ID        int       `json:"id"`
	Name      string    `json:"name"`
	Artist    string    `json:"artist"`
	Album     string    `json:"album"`
	Genre     string    `json:"genre"`
	CreatedAt time.Time `json:"created_at"`
}
