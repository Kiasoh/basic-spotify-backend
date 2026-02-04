package services

import (
	"context"
	"log"

	"github.com/kiasoh/basic-spotify-backend/models"
	"github.com/kiasoh/basic-spotify-backend/repository"
)

type SongService struct {
	Repo repository.SongRepository
}

func NewSongService(repo repository.SongRepository) *SongService {
	return &SongService{Repo: repo}
}

func (s *SongService) CreateSong(ctx context.Context, name, artist, album, genre string) (*models.Song, error) {
	log.Printf("Service: Attempting to create song: %s by %s", name, artist)
	song := &models.Song{
		Name:   name,
		Artist: artist,
		Album:  album,
		Genre:  genre,
	}

	id, err := s.Repo.CreateSong(ctx, song)
	if err != nil {
		log.Printf("Service: Error creating song in repository: %v", err)
		return nil, err
	}
	song.ID = id

	log.Printf("Service: Successfully created song with ID: %d", id)
	return song, nil
}

func (s *SongService) GetSong(ctx context.Context, id int) (*models.Song, error) {
	log.Printf("Service: Attempting to get song with ID: %d", id)
	song, err := s.Repo.GetSongByID(ctx, id)
	if err != nil {
		log.Printf("Service: Error getting song with ID %d: %v", id, err)
		return nil, err
	}
	return song, nil
}

func (s *SongService) ListSongs(ctx context.Context) ([]models.Song, error) {
	log.Println("Service: Attempting to list all songs")
	songs, err := s.Repo.ListSongs(ctx)
	if err != nil {
		log.Printf("Service: Error listing songs: %v", err)
		return nil, err
	}
	return songs, nil
}
