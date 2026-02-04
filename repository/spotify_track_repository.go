package repository

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/kiasoh/basic-spotify-backend/models"
)

type SpotifyTrackRepository interface {
	GetByTrackID(ctx context.Context, trackID string) (*models.SpotifyTrack, error)
	List(ctx context.Context, limit int, offset int, sortBy string) ([]models.SpotifyTrack, error)
	Search(ctx context.Context, query string, searchField string, limit int, offset int) ([]models.SpotifyTrack, error)
}

type spotifyTrackRepository struct {
	db *pgxpool.Pool
}

func NewSpotifyTrackRepository(db *pgxpool.Pool) SpotifyTrackRepository {
	return &spotifyTrackRepository{db: db}
}

func (r *spotifyTrackRepository) GetByTrackID(ctx context.Context, trackID string) (*models.SpotifyTrack, error) {
	query := `SELECT track_id, artists, album_name, track_name, popularity, duration_ms, explicit, danceability, energy, key, loudness, mode, speechiness, acousticness, instrumentalness, liveness, valence, tempo, time_signature, track_genre FROM spotify_tracks WHERE track_id = $1`
	track := &models.SpotifyTrack{}
	err := r.db.QueryRow(ctx, query, trackID).Scan(
		&track.TrackID, &track.Artists, &track.AlbumName, &track.TrackName, &track.Popularity, &track.DurationMs, &track.Explicit, &track.Danceability, &track.Energy, &track.Key, &track.Loudness, &track.Mode, &track.Speechiness, &track.Acousticness, &track.Instrumentalness, &track.Liveness, &track.Valence, &track.Tempo, &track.TimeSignature, &track.TrackGenre,
	)
	if err != nil {
		return nil, err
	}
	return track, nil
}

func (r *spotifyTrackRepository) List(ctx context.Context, limit int, offset int, sortBy string) ([]models.SpotifyTrack, error) {
	// Default sort by popularity if not specified or invalid
	if sortBy == "" {
		sortBy = "popularity"
	}

	// Validate sortBy field to prevent SQL injection
	// This list should ideally be dynamic or more comprehensive
	validSortFields := map[string]bool{
		"track_id": true, "artists": true, "album_name": true, "track_name": true,
		"popularity": true, "duration_ms": true, "explicit": true, "danceability": true,
		"energy": true, "key": true, "loudness": true, "mode": true,
		"speechiness": true, "acousticness": true, "instrumentalness": true,
		"liveness": true, "valence": true, "tempo": true, "time_signature": true,
		"track_genre": true,
	}

	if !validSortFields[sortBy] {
		return nil, fmt.Errorf("invalid sort field: %s", sortBy)
	}

	query := fmt.Sprintf(`SELECT track_id, artists, album_name, track_name, popularity, duration_ms, explicit, danceability, energy, key, loudness, mode, speechiness, acousticness, instrumentalness, liveness, valence, tempo, time_signature, track_genre FROM spotify_tracks ORDER BY %s DESC LIMIT $1 OFFSET $2`, sortBy)
	rows, err := r.db.Query(ctx, query, limit, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var tracks []models.SpotifyTrack
	for rows.Next() {
		var track models.SpotifyTrack
		if err := rows.Scan(
			&track.TrackID, &track.Artists, &track.AlbumName, &track.TrackName, &track.Popularity, &track.DurationMs, &track.Explicit, &track.Danceability, &track.Energy, &track.Key, &track.Loudness, &track.Mode, &track.Speechiness, &track.Acousticness, &track.Instrumentalness, &track.Liveness, &track.Valence, &track.Tempo, &track.TimeSignature, &track.TrackGenre,
		); err != nil {
			return nil, err
		}
		tracks = append(tracks, track)
	}
	return tracks, nil
}

func (r *spotifyTrackRepository) Search(ctx context.Context, query string, searchField string, limit int, offset int) ([]models.SpotifyTrack, error) {
	// Basic validation for searchField to prevent SQL injection
	switch searchField {
	case "track_name", "artists":
		// Valid fields
	default:
		return nil, fmt.Errorf("invalid search field: %s", searchField)
	}

	sqlQuery := fmt.Sprintf(`SELECT track_id, artists, album_name, track_name, popularity, duration_ms, explicit, danceability, energy, key, loudness, mode, speechiness, acousticness, instrumentalness, liveness, valence, tempo, time_signature, track_genre FROM spotify_tracks WHERE %s ILIKE '%%' || $1 || '%%' ORDER BY popularity DESC LIMIT $2 OFFSET $3`, searchField)
	
	rows, err := r.db.Query(ctx, sqlQuery, query, limit, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var tracks []models.SpotifyTrack
	for rows.Next() {
		var track models.SpotifyTrack
		if err := rows.Scan(
			&track.TrackID, &track.Artists, &track.AlbumName, &track.TrackName, &track.Popularity, &track.DurationMs, &track.Explicit, &track.Danceability, &track.Energy, &track.Key, &track.Loudness, &track.Mode, &track.Speechiness, &track.Acousticness, &track.Instrumentalness, &track.Liveness, &track.Valence, &track.Tempo, &track.TimeSignature, &track.TrackGenre,
		); err != nil {
			return nil, err
		}
		tracks = append(tracks, track)
	}
	return tracks, nil
}
