package services

import (
	"context"
	"errors"
	"fmt"
	"log"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/kiasoh/basic-spotify-backend/models"
	"github.com/kiasoh/basic-spotify-backend/repository"
)

type UserService struct {
	DB           *pgxpool.Pool
	UserRepo     repository.UserRepository
	PlaylistRepo repository.PlaylistRepository
}

func NewUserService(db *pgxpool.Pool, userRepo repository.UserRepository, playlistRepo repository.PlaylistRepository) *UserService {
	return &UserService{
		DB:           db,
		UserRepo:     userRepo,
		PlaylistRepo: playlistRepo,
	}
}

// RegisterUser handles the business logic for creating a new user and their default playlist in a transaction.
func (s *UserService) RegisterUser(ctx context.Context, username string, plaintextPassword string) (*models.User, error) {
	log.Printf("Attempting to register user: %s", username)

	// Validate password and check for existing user (outside the transaction)
	if valid, err := models.IsValidPassword(plaintextPassword); !valid {
		log.Printf("Validation error for user %s: %v", username, err)
		return nil, err
	}
	_, err := s.UserRepo.GetUserByUsername(ctx, username)
	if !errors.Is(err, pgx.ErrNoRows) {
		if err == nil {
			log.Printf("Registration failed for %s: user already exists", username)
			return nil, errors.New("user with this username already exists")
		}
		log.Printf("Error checking for existing user %s: %v", username, err)
		return nil, err
	}

	// Hash the password
	var password models.Password
	if err := password.Set(plaintextPassword); err != nil {
		log.Printf("Error hashing password for %s: %v", username, err)
		return nil, err
	}

	// --- Start Transaction ---
	tx, err := s.DB.Begin(ctx)
	if err != nil {
		log.Printf("Failed to begin transaction: %v", err)
		return nil, err
	}
	defer tx.Rollback(ctx)

	// 1. Create the user (without recomm_plylist_id)
	user := &models.User{
		Username: username,
		Password: password,
	}
	userID, err := s.UserRepo.CreateUserInTx(ctx, tx, user)
	if err != nil {
		log.Printf("Error creating user %s in transaction: %v", username, err)
		return nil, err
	}
	user.ID = userID

	// 2. Create the default playlist, using the new userID as the OwnerID
	defaultPlaylist := &models.Playlist{
		Name:    fmt.Sprintf("%s's Recommendations", username),
		OwnerID: userID,
	}
	playlistID, err := s.PlaylistRepo.CreatePlaylistInTx(ctx, tx, defaultPlaylist)
	if err != nil {
		log.Printf("Failed to create default playlist for user %s: %v", username, err)
		return nil, err
	}
	user.RecommPlaylistID = playlistID

	// 3. Update the user with the new playlist ID
	err = s.UserRepo.UpdateRecommPlaylistIDInTx(ctx, tx, userID, playlistID)
	if err != nil {
		log.Printf("Failed to update user %d with recomm playlist ID %d: %v", userID, playlistID, err)
		return nil, err
	}

	// --- Commit Transaction ---
	if err := tx.Commit(ctx); err != nil {
		log.Printf("Failed to commit transaction for user %s: %v", username, err)
		return nil, err
	}

	log.Printf("Successfully registered user %s with ID: %d and default playlist ID: %d", username, userID, playlistID)
	return user, nil
}
