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
	AddSongToPlaylist(ctx context.Context, playlistID int, songID int) error
	RemoveSongFromPlaylist(ctx context.Context, playlistID int, songID int) error
	GetSongsInPlaylist(ctx context.Context, playlistID int) ([]models.Song, error)
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

func (r *playlistRepository) AddSongToPlaylist(ctx context.Context, playlistID int, songID int) error {
	query := `INSERT INTO songs_playlists (playlist_id, song_id) VALUES ($1, $2)`
	_, err := r.db.Exec(ctx, query, playlistID, songID)
	return err
}

func (r *playlistRepository) RemoveSongFromPlaylist(ctx context.Context, playlistID int, songID int) error {
	query := `DELETE FROM songs_playlists WHERE playlist_id = $1 AND song_id = $2`
	_, err := r.db.Exec(ctx, query, playlistID, songID)
	return err
}

func (r *playlistRepository) GetSongsInPlaylist(ctx context.Context, playlistID int) ([]models.Song, error) {
	query := `
		SELECT s.id, s.name, s.artist, s.album, s.genre, s.created_at 
		FROM songs s 
		JOIN songs_playlists sp ON s.id = sp.song_id 
		WHERE sp.playlist_id = $1`
	rows, err := r.db.Query(ctx, query, playlistID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var songs []models.Song
	for rows.Next() {
		var song models.Song
		if err := rows.Scan(&song.ID, &song.Name, &song.Artist, &song.Album, &song.Genre, &song.CreatedAt); err != nil {
			return nil, err
		}
		songs = append(songs, song)
	}
	return songs, nil
}
