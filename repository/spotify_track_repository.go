package repository

import (
	"context"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/kiasoh/basic-spotify-backend/models"
)

type SpotifyTrackRepository interface {
	GetByTrackID(ctx context.Context, trackID string) (*models.SpotifyTrack, error)
	List(ctx context.Context, limit int, offset int) ([]models.SpotifyTrack, error)
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

func (r *spotifyTrackRepository) List(ctx context.Context, limit int, offset int) ([]models.SpotifyTrack, error) {
	query := `SELECT track_id, artists, album_name, track_name, popularity, duration_ms, explicit, danceability, energy, key, loudness, mode, speechiness, acousticness, instrumentalness, liveness, valence, tempo, time_signature, track_genre FROM spotify_tracks ORDER BY popularity DESC LIMIT $1 OFFSET $2`
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
