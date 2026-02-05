package repository

import (
	"context"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/kiasoh/basic-spotify-backend/models"
)

type PlaylistRepository interface {
	CreatePlaylist(ctx context.Context, playlist *models.Playlist) (int, error)
	CreatePlaylistInTx(ctx context.Context, tx pgx.Tx, playlist *models.Playlist) (int, error)
	GetPlaylistByID(ctx context.Context, id int) (*models.Playlist, error)
	UpdatePlaylist(ctx context.Context, playlist *models.Playlist) error
	DeletePlaylist(ctx context.Context, id int) error
	ListPlaylistsByOwner(ctx context.Context, ownerID int) ([]models.Playlist, error)
	AddTrackToPlaylist(ctx context.Context, playlistID int, trackID string) error
	RemoveTrackFromPlaylist(ctx context.Context, playlistID int, trackID string) error
	GetTracksInPlaylist(ctx context.Context, playlistID int) ([]models.SpotifyTrack, error)
	GetTrackInPlaylist(ctx context.Context, playlistID int, trackID string) (*models.SpotifyTrack, error)
}

type playlistRepository struct {
	db *pgxpool.Pool
}

func NewPlaylistRepository(db *pgxpool.Pool) PlaylistRepository {
	return &playlistRepository{db: db}
}

func (r *playlistRepository) CreatePlaylist(ctx context.Context, playlist *models.Playlist) (int, error) {
	query := `INSERT INTO playlists (name, owner_id, modifyable, description) VALUES ($1, $2, $3, $4) RETURNING id`
	var id int
	err := r.db.QueryRow(ctx, query, playlist.Name, playlist.OwnerID, playlist.Modifyable, playlist.Description).Scan(&id)
	return id, err
}

func (r *playlistRepository) CreatePlaylistInTx(ctx context.Context, tx pgx.Tx, playlist *models.Playlist) (int, error) {
	query := `INSERT INTO playlists (name, owner_id, modifyable, description) VALUES ($1, $2, $3, $4) RETURNING id`
	var id int
	err := tx.QueryRow(ctx, query, playlist.Name, playlist.OwnerID, playlist.Modifyable, playlist.Description).Scan(&id)
	return id, err
}

func (r *playlistRepository) GetPlaylistByID(ctx context.Context, id int) (*models.Playlist, error) {
	query := `SELECT id, name, description, owner_id, modifyable, created_at FROM playlists WHERE id = $1`
	playlist := &models.Playlist{}
	err := r.db.QueryRow(ctx, query, id).Scan(&playlist.ID, &playlist.Name, &playlist.Description, &playlist.OwnerID, &playlist.Modifyable, &playlist.CreatedAt)
	if err != nil {
		return nil, err
	}
	return playlist, nil
}

func (r *playlistRepository) UpdatePlaylist(ctx context.Context, playlist *models.Playlist) error {
	query := `UPDATE playlists SET name = $1, modifyable = $2, description = $3 WHERE id = $4`
	_, err := r.db.Exec(ctx, query, playlist.Name, playlist.Modifyable, playlist.Description, playlist.ID)
	return err
}

func (r *playlistRepository) DeletePlaylist(ctx context.Context, id int) error {
	tx, err := r.db.Begin(ctx)
	if err != nil {
		return err
	}
	defer tx.Rollback(ctx)

	_, err = tx.Exec(ctx, "DELETE FROM songs_playlists WHERE playlist_id = $1", id)
	if err != nil {
		return err
	}

	_, err = tx.Exec(ctx, "DELETE FROM playlists WHERE id = $1", id)
	if err != nil {
		return err
	}

	return tx.Commit(ctx)
}

func (r *playlistRepository) ListPlaylistsByOwner(ctx context.Context, ownerID int) ([]models.Playlist, error) {
	query := `SELECT id, name, description, owner_id, modifyable, created_at FROM playlists WHERE owner_id = $1`
	rows, err := r.db.Query(ctx, query, ownerID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var playlists []models.Playlist
	for rows.Next() {
		var playlist models.Playlist
		if err := rows.Scan(&playlist.ID, &playlist.Name, &playlist.Description, &playlist.OwnerID, &playlist.Modifyable, &playlist.CreatedAt); err != nil {
			return nil, err
		}
		playlists = append(playlists, playlist)
	}
	return playlists, nil
}

func (r *playlistRepository) AddTrackToPlaylist(ctx context.Context, playlistID int, trackID string) error {
	query := `INSERT INTO songs_playlists (playlist_id, track_id) VALUES ($1, $2)`
	_, err := r.db.Exec(ctx, query, playlistID, trackID)
	return err
}

func (r *playlistRepository) RemoveTrackFromPlaylist(ctx context.Context, playlistID int, trackID string) error {
	query := `DELETE FROM songs_playlists WHERE playlist_id = $1 AND track_id = $2`
	_, err := r.db.Exec(ctx, query, playlistID, trackID)
	return err
}

func (r *playlistRepository) GetTrackInPlaylist(ctx context.Context, playlistID int, trackID string) (*models.SpotifyTrack, error) {
	query := `
		SELECT t.track_id, t.artists, t.album_name, t.track_name, t.popularity, t.duration_ms, t.explicit, t.danceability, t.energy, t.key, t.loudness, t.mode, t.speechiness, t.acousticness, t.instrumentalness, t.liveness, t.valence, t.tempo, t.time_signature, t.track_genre
		FROM spotify_tracks 
		WHERE playlist_id = $1 and track_id = $2`

	var track models.SpotifyTrack
	err := r.db.QueryRow(ctx, query, playlistID, trackID).Scan(
		&track.TrackID, &track.Artists, &track.AlbumName, &track.TrackName, &track.Popularity, &track.DurationMs, &track.Explicit, &track.Danceability, &track.Energy, &track.Key, &track.Loudness, &track.Mode, &track.Speechiness, &track.Acousticness, &track.Instrumentalness, &track.Liveness, &track.Valence, &track.Tempo, &track.TimeSignature, &track.TrackGenre)
	if err != nil {
		return nil, err
	}
	return &track, nil
}

func (r *playlistRepository) GetTracksInPlaylist(ctx context.Context, playlistID int) ([]models.SpotifyTrack, error) {
	query := `
		SELECT t.track_id, t.artists, t.album_name, t.track_name, t.popularity, t.duration_ms, t.explicit, t.danceability, t.energy, t.key, t.loudness, t.mode, t.speechiness, t.acousticness, t.instrumentalness, t.liveness, t.valence, t.tempo, t.time_signature, t.track_genre
		FROM spotify_tracks t
		JOIN songs_playlists sp ON t.track_id = sp.track_id
		WHERE sp.playlist_id = $1`
	rows, err := r.db.Query(ctx, query, playlistID)
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
