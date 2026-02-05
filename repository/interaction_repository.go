package repository

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/kiasoh/basic-spotify-backend/models"
)

type InteractionRepository interface {
	CreateInteraction(ctx context.Context, interaction *models.Interaction) error
	GetInteractionsByUser(ctx context.Context, userID int) ([]models.Interaction, error)
	GetInteractionsForTrack(ctx context.Context, trackID string) ([]models.Interaction, error)
	GetLatestInteractionsForUserTracks(ctx context.Context, userID int, trackIDs []string) (map[string]string, error)
}

type interactionRepository struct {
	db *pgxpool.Pool
}

func NewInteractionRepository(db *pgxpool.Pool) InteractionRepository {
	return &interactionRepository{db: db}
}

func (r *interactionRepository) CreateInteraction(ctx context.Context, interaction *models.Interaction) error {
	query := `INSERT INTO interactions (user_id, track_id, type) VALUES ($1, $2, $3)`
	_, err := r.db.Exec(ctx, query, interaction.UserID, interaction.TrackID, interaction.Type)
	return err
}

func (r *interactionRepository) GetInteractionsByUser(ctx context.Context, userID int) ([]models.Interaction, error) {
	query := `SELECT user_id, track_id, type, created_at FROM interactions WHERE user_id = $1`
	rows, err := r.db.Query(ctx, query, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var interactions []models.Interaction
	for rows.Next() {
		var i models.Interaction
		if err := rows.Scan(&i.UserID, &i.TrackID, &i.Type, &i.CreatedAt); err != nil {
			return nil, err
		}
		interactions = append(interactions, i)
	}
	return interactions, nil
}

func (r *interactionRepository) GetInteractionsForTrack(ctx context.Context, trackID string) ([]models.Interaction, error) {
	query := `SELECT user_id, track_id, type, created_at FROM interactions WHERE track_id = $1`
	rows, err := r.db.Query(ctx, query, trackID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var interactions []models.Interaction
	for rows.Next() {
		var i models.Interaction
		if err := rows.Scan(&i.UserID, &i.TrackID, &i.Type, &i.CreatedAt); err != nil {
			return nil, err
		}
		interactions = append(interactions, i)
	}
	return interactions, nil
}

func (r *interactionRepository) GetLatestInteractionsForUserTracks(ctx context.Context, userID int, trackIDs []string) (map[string]string, error) {
	if len(trackIDs) == 0 {
		return make(map[string]string), nil
	}

	query := `
		SELECT DISTINCT ON (track_id)
			track_id, type
		FROM
			interactions
		WHERE
			user_id = $1 AND track_id = ANY($2)
		ORDER BY
			track_id, created_at DESC
	`
	rows, err := r.db.Query(ctx, query, userID, trackIDs)
	if err != nil {
		return nil, fmt.Errorf("failed to get latest interactions: %w", err)
	}
	defer rows.Close()

	interactionMap := make(map[string]string)
	for rows.Next() {
		var trackID, interactionType string
		if err := rows.Scan(&trackID, &interactionType); err != nil {
			return nil, fmt.Errorf("failed to scan interaction row: %w", err)
		}
		interactionMap[trackID] = interactionType
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error during rows iteration: %w", err)
	}

	return interactionMap, nil
}
