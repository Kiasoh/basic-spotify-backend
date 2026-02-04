package models

type SpotifyTrack struct {
	TrackID          string  `json:"track_id"`
	Artists          string  `json:"artists"`
	AlbumName        string  `json:"album_name"`
	TrackName        string  `json:"track_name"`
	Popularity       int64   `json:"popularity"`
	DurationMs       int64   `json:"duration_ms"`
	Explicit         bool    `json:"explicit"`
	Danceability     float64 `json:"danceability"`
	Energy           float64 `json:"energy"`
	Key              int64   `json:"key"`
	Loudness         float64 `json:"loudness"`
	Mode             int64   `json:"mode"`
	Speechiness      float64 `json:"speechiness"`
	Acousticness     float64 `json:"acousticness"`
	Instrumentalness float64 `json:"instrumentalness"`
	Liveness         float64 `json:"liveness"`
	Valence          float64 `json:"valence"`
	Tempo            float64 `json:"tempo"`
	TimeSignature    int64   `json:"time_signature"`
	TrackGenre       string  `json:"track_genre"`
}
