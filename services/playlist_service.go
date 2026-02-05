package services

import (
	"context"
	"errors"
	"log"

	"github.com/kiasoh/basic-spotify-backend/models"
	"github.com/kiasoh/basic-spotify-backend/repository"
)

type PlaylistService struct {
	Repo             repository.PlaylistRepository
	InteractionService *InteractionService // Add this field
}

func NewPlaylistService(repo repository.PlaylistRepository, interactionService *InteractionService) *PlaylistService {
	return &PlaylistService{Repo: repo, InteractionService: interactionService}
}

func (s *PlaylistService) CreatePlaylist(ctx context.Context, ownerID int, name string, description *string) (*models.Playlist, error) {
	log.Printf("Service: User %d attempting to create playlist '%s'", ownerID, name)
	playlist := &models.Playlist{
		Name:        name,
		OwnerID:     ownerID,
		Description: description,
		Modifyable:  true, // User-created playlists are always modifiable
	}

	id, err := s.Repo.CreatePlaylist(ctx, playlist)
	if err != nil {
		log.Printf("Service: Error creating playlist: %v", err)
		return nil, err
	}
	playlist.ID = id
	log.Printf("Service: Successfully created playlist with ID %d for user %d", id, ownerID)
	return playlist, nil
}

func (s *PlaylistService) ListUserPlaylists(ctx context.Context, ownerID int) ([]models.Playlist, error) {
	log.Printf("Service: User %d attempting to list their playlists", ownerID)
	return s.Repo.ListPlaylistsByOwner(ctx, ownerID)
}

func (s *PlaylistService) AddTrackToPlaylist(ctx context.Context, userID, playlistID int, trackID string) error {
	log.Printf("Service: User %d attempting to add track %s to playlist %d", userID, trackID, playlistID)

	playlist, err := s.Repo.GetPlaylistByID(ctx, playlistID)
	if err != nil {
		log.Printf("Service: Error checking playlist: could not get playlist %d. Error: %v", playlistID, err)
		return errors.New("playlist not found")
	}
	if playlist.OwnerID != userID {
		log.Printf("Service: User %d does not own playlist %d", userID, playlistID)
		return errors.New("forbidden: you do not own this playlist")
	}
	if !playlist.Modifyable {
		log.Printf("Service: User %d cannot modify unmodifiable playlist %d", userID, playlistID)
		return errors.New("forbidden: this playlist is not modifiable")
	}

	return s.Repo.AddTrackToPlaylist(ctx, playlistID, trackID)
}

func (s *PlaylistService) RemoveTrackFromPlaylist(ctx context.Context, userID, playlistID int, trackID string) error {
	log.Printf("Service: User %d attempting to remove track %s from playlist %d", userID, trackID, playlistID)

	playlist, err := s.Repo.GetPlaylistByID(ctx, playlistID)
	if err != nil {
		log.Printf("Service: Error checking playlist: could not get playlist %d. Error: %v", playlistID, err)
		return errors.New("playlist not found")
	}
	if playlist.OwnerID != userID {
		log.Printf("Service: User %d does not own playlist %d", userID, playlistID)
		return errors.New("forbidden: you do not own this playlist")
	}
	if !playlist.Modifyable {
		log.Printf("Service: User %d cannot modify unmodifiable playlist %d", userID, playlistID)
		return errors.New("forbidden: this playlist is not modifiable")
	}

	return s.Repo.RemoveTrackFromPlaylist(ctx, playlistID, trackID)
}

func (s *PlaylistService) UpdatePlaylistDetails(ctx context.Context, userID int, playlistID int, newName string, newDescription *string) (*models.Playlist, error) {
	log.Printf("Service: User %d attempting to update details for playlist %d", userID, playlistID)

	playlist, err := s.Repo.GetPlaylistByID(ctx, playlistID)
	if err != nil {
		log.Printf("Service: Error checking playlist: could not get playlist %d. Error: %v", playlistID, err)
		return nil, errors.New("playlist not found")
	}
	if playlist.OwnerID != userID {
		log.Printf("Service: User %d does not own playlist %d", userID, playlistID)
		return nil, errors.New("forbidden: you do not own this playlist")
	}
	if !playlist.Modifyable {
		log.Printf("Service: User %d cannot modify unmodifiable playlist %d", userID, playlistID)
		return nil, errors.New("forbidden: this playlist is not modifiable")
	}

	// Update fields
	playlist.Name = newName
	playlist.Description = newDescription

	err = s.Repo.UpdatePlaylist(ctx, playlist)
	if err != nil {
		log.Printf("Service: Failed to update playlist %d: %v", playlistID, err)
		return nil, err
	}

	return playlist, nil
}

func (s *PlaylistService) GetTracksInPlaylist(ctx context.Context, playlistID int) ([]models.SpotifyTrack, error) {
	log.Printf("Service: Attempting to get tracks for playlist %d", playlistID)
	return s.Repo.GetTracksInPlaylist(ctx, playlistID)
}
