package services

import (
	"context"
	"errors"
	"log"
	"strconv"

	"github.com/kiasoh/basic-spotify-backend/models"
	"github.com/kiasoh/basic-spotify-backend/repository"
	"github.com/segmentio/kafka-go"
)

type InteractionService struct {
	Repo        repository.InteractionRepository
	TrackRepo   repository.SpotifyTrackRepository
	UserRepo    repository.UserRepository
	KafkaWriter *kafka.Writer
}

func NewInteractionService(repo repository.InteractionRepository, kafkaWriter *kafka.Writer, trackRepo repository.SpotifyTrackRepository, userRepo repository.UserRepository) *InteractionService {
	return &InteractionService{
		Repo:        repo,
		TrackRepo:   trackRepo,
		UserRepo:    userRepo,
		KafkaWriter: kafkaWriter,
	}
}

func (s *InteractionService) convertTrackToVector(track *models.SpotifyTrack, multiplier float64) ([]float64, error) {
	return []float64{track.Danceability * multiplier, track.Energy * multiplier, track.Loudness * multiplier, track.Speechiness * multiplier, track.Acousticness * multiplier, track.Instrumentalness * multiplier, track.Liveness * multiplier, track.Valence * multiplier, track.Tempo * multiplier}, nil
}

const alpha = 0.45

func (s *InteractionService) HandleInteraction(ctx context.Context, userID int, trackID string, interactionType string) error {
	var weight float64
	switch interactionType {
	case "like":
		weight = 3.0
	case "unlike":
		weight = -2.5
	case "dislike":
		weight = -4.0
	case "undislike":
		weight = 4.5
	case "skip":
		weight = -1.0
	case "play":
		weight = 1.0
	case "add_to_playlist":
		weight = 5.0
	case "remove_from_playlist":
		weight = -3.0
	default:
		return errors.New("invalid interaction type")
	}
	track, err := s.TrackRepo.GetByTrackID(ctx, trackID)
	if err != nil {
		log.Println("Service: Track not found", err)
		return err
	}
	user, err := s.UserRepo.GetUserByID(ctx, userID)
	if err != nil {
		log.Println("Service: User not found", err)
		return err
	}

	trackVector, err := s.convertTrackToVector(track, weight)
	if err != nil {
		log.Println("Service: Error converting track to vector", err)
		return err
	}

	for i := 0; i < len(user.AvgInterest); i++ {
		user.AvgInterest[i] = alpha*user.AvgInterest[i] + (1-alpha)*trackVector[i]
	}

	err = s.UserRepo.UpdateUser(ctx, user)
	if err != nil {
		log.Println("Service: Error updating user interest", err)
		return err
	}

	return nil
}

func (s *InteractionService) CreateInteraction(ctx context.Context, userID int, trackID string, interactionType string) error {
	log.Printf("Service: User %d creating interaction of type '%s' for track %s", userID, interactionType, trackID)

	err := s.HandleInteraction(ctx, userID, trackID, interactionType)
	if err != nil {
		log.Printf("Service: Error handling interaction in service: %v", err)
		return err
	}

	interaction := &models.Interaction{
		UserID:  userID,
		TrackID: trackID,
		Type:    interactionType,
	}

	err = s.Repo.CreateInteraction(ctx, interaction)
	if err != nil {
		log.Printf("Service: Error creating interaction in DB: %v", err)
		return err
	}

	log.Printf("Service: Successfully created interaction for user %d and track %s. Publishing to Kafka.", userID, trackID)

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

func (s *InteractionService) GetInteractionsForTrack(ctx context.Context, trackID string) ([]models.Interaction, error) {
	log.Printf("Service: Getting interactions for track %s", trackID)
	return s.Repo.GetInteractionsForTrack(ctx, trackID)
}

func (s *InteractionService) GetTrackInteractionStates(ctx context.Context, userID int, trackIDs []string) (map[string]models.TrackInteractionState, error) {
	interactionMap, err := s.Repo.GetLatestInteractionsForUserTracks(ctx, userID, trackIDs)
	if err != nil {
		return nil, err
	}

	trackStates := make(map[string]models.TrackInteractionState)
	for _, trackID := range trackIDs {
		interactionType, found := interactionMap[trackID]
		if !found {
			trackStates[trackID] = models.TrackStateNeutral
			continue
		}

		switch interactionType {
		case "like":
			trackStates[trackID] = models.TrackStateLiked
		case "dislike":
			trackStates[trackID] = models.TrackStateDisliked
		default:
			trackStates[trackID] = models.TrackStateNeutral
		}
	}
	return trackStates, nil
}
