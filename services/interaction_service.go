package services

import (
	"context"
	"log"
	"strconv"

	"github.com/kiasoh/basic-spotify-backend/models"
	"github.com/kiasoh/basic-spotify-backend/repository"
	"github.com/segmentio/kafka-go"
)

type InteractionService struct {
	Repo        repository.InteractionRepository
	KafkaWriter *kafka.Writer
}

func NewInteractionService(repo repository.InteractionRepository, kafkaWriter *kafka.Writer) *InteractionService {
	return &InteractionService{
		Repo:        repo,
		KafkaWriter: kafkaWriter,
	}
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
		log.Printf("Service: Error creating interaction in DB: %v", err)
		return err
	}

	log.Printf("Service: Successfully created interaction for user %d and song %d. Publishing to Kafka.", userID, songID)

	// Publish event to Kafka
	msg := kafka.Message{
		Value: []byte(strconv.Itoa(userID)),
	}
	err = s.KafkaWriter.WriteMessages(ctx, msg)
	if err != nil {
		// Log the error but don't return it to the client,
		// as the primary operation (saving the interaction) was successful.
		log.Printf("Service: Failed to write message to Kafka: %v", err)
	}

	return nil
}

func (s *InteractionService) GetInteractionsForSong(ctx context.Context, songID int) ([]models.Interaction, error) {
	log.Printf("Service: Getting interactions for song %d", songID)
	return s.Repo.GetInteractionsForSong(ctx, songID)
}
