package repository

import (
	"context"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/kiasoh/basic-spotify-backend/models"
)

type SongRepository interface {
	CreateSong(ctx context.Context, song *models.Song) (int, error)
	GetSongByID(ctx context.Context, id int) (*models.Song, error)
	UpdateSong(ctx context.Context, song *models.Song) error
	DeleteSong(ctx context.Context, id int) error
	ListSongs(ctx context.Context) ([]models.Song, error)
}

type songRepository struct {
	db *pgxpool.Pool
}

func NewSongRepository(db *pgxpool.Pool) SongRepository {
	return &songRepository{db: db}
}

func (r *songRepository) CreateSong(ctx context.Context, song *models.Song) (int, error) {
	query := `INSERT INTO songs (name, artist, album, genre) VALUES ($1, $2, $3, $4) RETURNING id`
	var id int
	err := r.db.QueryRow(ctx, query, song.Name, song.Artist, song.Album, song.Genre).Scan(&id)
	return id, err
}

func (r *songRepository) GetSongByID(ctx context.Context, id int) (*models.Song, error) {
	query := `SELECT id, name, artist, album, genre, created_at FROM songs WHERE id = $1`
	song := &models.Song{}
	err := r.db.QueryRow(ctx, query, id).Scan(&song.ID, &song.Name, &song.Artist, &song.Album, &song.Genre, &song.CreatedAt)
	if err != nil {
		return nil, err
	}
	return song, nil
}

func (r *songRepository) UpdateSong(ctx context.Context, song *models.Song) error {
	query := `UPDATE songs SET name = $1, artist = $2, album = $3, genre = $4 WHERE id = $5`
	_, err := r.db.Exec(ctx, query, song.Name, song.Artist, song.Album, song.Genre, song.ID)
	return err
}

func (r *songRepository) DeleteSong(ctx context.Context, id int) error {
	query := `DELETE FROM songs WHERE id = $1`
	_, err := r.db.Exec(ctx, query, id)
	return err
}

func (r *songRepository) ListSongs(ctx context.Context) ([]models.Song, error) {
	query := `SELECT id, name, artist, album, genre, created_at FROM songs`
	rows, err := r.db.Query(ctx, query)
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
