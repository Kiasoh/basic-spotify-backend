package services

import (
	"context"
	"log"

	"github.com/kiasoh/basic-spotify-backend/models"
	"github.com/kiasoh/basic-spotify-backend/repository"
)

type SpotifyTrackService struct {
	Repo repository.SpotifyTrackRepository
}

func NewSpotifyTrackService(repo repository.SpotifyTrackRepository) *SpotifyTrackService {
	return &SpotifyTrackService{Repo: repo}
}

func (s *SpotifyTrackService) GetByTrackID(ctx context.Context, trackID string) (*models.SpotifyTrack, error) {
	log.Printf("Service: Attempting to get track with ID: %s", trackID)
	track, err := s.Repo.GetByTrackID(ctx, trackID)
	if err != nil {
		log.Printf("Service: Error getting track with ID %s: %v", trackID, err)
		return nil, err
	}
	return track, nil
}

func (s *SpotifyTrackService) List(ctx context.Context, limit int, offset int, sortBy string, order string) ([]models.SpotifyTrack, error) {
	log.Printf("Service: Attempting to list tracks with limit %d, offset %d, and sort by %s", limit, offset, sortBy)
	tracks, err := s.Repo.List(ctx, limit, offset, sortBy, order)
	if err != nil {
		log.Printf("Service: Error listing tracks: %v", err)
		return nil, err
	}
	return tracks, nil
}

func (s *SpotifyTrackService) Search(ctx context.Context, query string, searchField string, limit int, offset int) ([]models.SpotifyTrack, error) {
	log.Printf("Service: Attempting to search tracks for query '%s' in field '%s' with limit %d and offset %d", query, searchField, limit, offset)
	tracks, err := s.Repo.Search(ctx, query, searchField, limit, offset)
	if err != nil {
		log.Printf("Service: Error searching tracks: %v", err)
		return nil, err
	}
	return tracks, nil
}
