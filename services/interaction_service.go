package services

import (
	"context"
	"log"

	"github.com/kiasoh/basic-spotify-backend/models"
	"github.com/kiasoh/basic-spotify-backend/repository"
)

type InteractionService struct {
	Repo repository.InteractionRepository
}

func NewInteractionService(repo repository.InteractionRepository) *InteractionService {
	return &InteractionService{Repo: repo}
}

func (s *InteractionService) CreateInteraction(ctx context.Context, userID, songID int, interactionType string) error {
	log.Printf("Service: User %d creating interaction of type '%s' for song %d", userID, interactionType, songID)

	interaction := &models.Interaction{
		UserID: userID,
		SongID: songID,
		Type:   interactionType,
	}

	err := s.Repo.CreateInteraction(ctx, interaction)
	if err != nil {
		log.Printf("Service: Error creating interaction: %v", err)
		return err
	}

	log.Printf("Service: Successfully created interaction for user %d and song %d", userID, songID)
	return nil
}

func (s *InteractionService) GetInteractionsForSong(ctx context.Context, songID int) ([]models.Interaction, error) {
	log.Printf("Service: Getting interactions for song %d", songID)
	return s.Repo.GetInteractionsForSong(ctx, songID)
}
