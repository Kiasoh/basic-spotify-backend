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

func (s *SpotifyTrackService) List(ctx context.Context, limit int, offset int) ([]models.SpotifyTrack, error) {
	log.Printf("Service: Attempting to list tracks with limit %d and offset %d", limit, offset)
	tracks, err := s.Repo.List(ctx, limit, offset)
	if err != nil {
		log.Printf("Service: Error listing tracks: %v", err)
		return nil, err
	}
	return tracks, nil
}
