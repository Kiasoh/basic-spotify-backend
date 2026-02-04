package services

import (
	"context"
	"errors"
	"log"

	"github.com/kiasoh/basic-spotify-backend/models"
	"github.com/kiasoh/basic-spotify-backend/repository"
)

type PlaylistService struct {
	Repo repository.PlaylistRepository
}

func NewPlaylistService(repo repository.PlaylistRepository) *PlaylistService {
	return &PlaylistService{Repo: repo}
}

func (s *PlaylistService) CreatePlaylist(ctx context.Context, name string, ownerID int) (*models.Playlist, error) {
	log.Printf("Service: User %d attempting to create playlist '%s'", ownerID, name)
	playlist := &models.Playlist{
		Name:    name,
		OwnerID: ownerID,
		// Modifyable defaults to true in the DB
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

func (s *PlaylistService) AddSongToPlaylist(ctx context.Context, userID, playlistID, songID int) error {
	log.Printf("Service: User %d attempting to add song %d to playlist %d", userID, songID, playlistID)

	// Business rule: Verify the user owns the playlist
	playlist, err := s.Repo.GetPlaylistByID(ctx, playlistID)
	if err != nil {
		log.Printf("Service: Error checking playlist ownership: could not get playlist %d. Error: %v", playlistID, err)
		return errors.New("playlist not found")
	}
	if playlist.OwnerID != userID {
		log.Printf("Service: User %d does not own playlist %d", userID, playlistID)
		return errors.New("forbidden: you do not own this playlist")
	}

	return s.Repo.AddSongToPlaylist(ctx, playlistID, songID)
}

func (s *PlaylistService) RemoveSongFromPlaylist(ctx context.Context, userID, playlistID, songID int) error {
	log.Printf("Service: User %d attempting to remove song %d from playlist %d", userID, songID, playlistID)

	// Business rule: Verify the user owns the playlist
	playlist, err := s.Repo.GetPlaylistByID(ctx, playlistID)
	if err != nil {
		log.Printf("Service: Error checking playlist ownership: could not get playlist %d. Error: %v", playlistID, err)
		return errors.New("playlist not found")
	}
	if playlist.OwnerID != userID {
		log.Printf("Service: User %d does not own playlist %d", userID, playlistID)
		return errors.New("forbidden: you do not own this playlist")
	}

	return s.Repo.RemoveSongFromPlaylist(ctx, playlistID, songID)
}

func (s *PlaylistService) GetSongsInPlaylist(ctx context.Context, playlistID int) ([]models.Song, error) {
	log.Printf("Service: Attempting to get songs for playlist %d", playlistID)
	// For this, we are not checking ownership, assuming playlists can be viewed publicly.
	// This could be changed by adding a userID check similar to Add/Remove.
	return s.Repo.GetSongsInPlaylist(ctx, playlistID)
}
