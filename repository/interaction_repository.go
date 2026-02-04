package repository

import (
	"context"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/kiasoh/basic-spotify-backend/models"
)

type InteractionRepository interface {
	CreateInteraction(ctx context.Context, interaction *models.Interaction) error
	GetInteractionsByUser(ctx context.Context, userID int) ([]models.Interaction, error)
    GetInteractionsForSong(ctx context.Context, songID int) ([]models.Interaction, error)
}

type interactionRepository struct {
	db *pgxpool.Pool
}

func NewInteractionRepository(db *pgxpool.Pool) InteractionRepository {
	return &interactionRepository{db: db}
}

func (r *interactionRepository) CreateInteraction(ctx context.Context, interaction *models.Interaction) error {
	query := `INSERT INTO interactions (user_id, song_id, type) VALUES ($1, $2, $3)`
	_, err := r.db.Exec(ctx, query, interaction.UserID, interaction.SongID, interaction.Type)
	return err
}

func (r *interactionRepository) GetInteractionsByUser(ctx context.Context, userID int) ([]models.Interaction, error) {
	query := `SELECT user_id, song_id, type, created_at FROM interactions WHERE user_id = $1`
	rows, err := r.db.Query(ctx, query, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var interactions []models.Interaction
	for rows.Next() {
		var i models.Interaction
		if err := rows.Scan(&i.UserID, &i.SongID, &i.Type, &i.CreatedAt); err != nil {
			return nil, err
		}
		interactions = append(interactions, i)
	}
	return interactions, nil
}

func (r *interactionRepository) GetInteractionsForSong(ctx context.Context, songID int) ([]models.Interaction, error) {
	query := `SELECT user_id, song_id, type, created_at FROM interactions WHERE song_id = $1`
	rows, err := r.db.Query(ctx, query, songID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var interactions []models.Interaction
	for rows.Next() {
		var i models.Interaction
		if err := rows.Scan(&i.UserID, &i.SongID, &i.Type, &i.CreatedAt); err != nil {
			return nil, err
		}
		interactions = append(interactions, i)
	}
	return interactions, nil
}
