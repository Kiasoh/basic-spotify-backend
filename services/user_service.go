package services

import (
	"context"
	"errors"
	"log"

	"github.com/jackc/pgx/v5"
	"github.com/kiasoh/basic-spotify-backend/models"
	"github.com/kiasoh/basic-spotify-backend/repository"
)

type UserService struct {
	Repo repository.UserRepository
}

func NewUserService(repo repository.UserRepository) *UserService {
	return &UserService{Repo: repo}
}

// RegisterUser handles the business logic for creating a new user.
func (s *UserService) RegisterUser(ctx context.Context, username string, plaintextPassword string) (*models.User, error) {
	log.Printf("Attempting to register user: %s", username)

	// Validate password
	if valid, err := models.IsValidPassword(plaintextPassword); !valid {
		log.Printf("Validation error for user %s: %v", username, err)
		return nil, err
	}

	// Check if user already exists
	_, err := s.Repo.GetUserByUsername(ctx, username)
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

	user := &models.User{
		Username: username,
		Password: password,
	}

	// Create user in the repository
	id, err := s.Repo.CreateUser(ctx, user)
	if err != nil {
		log.Printf("Error creating user %s in repository: %v", username, err)
		return nil, err
	}
	user.ID = id
	user.Password.Plaintext = &plaintextPassword // Keep plaintext for potential immediate use, though it won't be stored

	log.Printf("Successfully registered user %s with ID: %d", username, id)
	return user, nil
}
