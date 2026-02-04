package services

import (
	"context"
	"errors"
	"log"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/kiasoh/basic-spotify-backend/repository"
)

// TODO: Move this to a secure location like environment variables
var jwtSecret = []byte("supersecretkey")

type AuthService struct {
	UserRepo repository.UserRepository
}

func NewAuthService(userRepo repository.UserRepository) *AuthService {
	return &AuthService{UserRepo: userRepo}
}

// Login validates user credentials and returns a JWT token if they are correct.
func (s *AuthService) Login(ctx context.Context, username, plaintextPassword string) (string, error) {
	log.Printf("Attempting login for user: %s", username)

	user, err := s.UserRepo.GetUserByUsername(ctx, username)
	if err != nil {
		log.Printf("Login failed for %s: user not found. Error: %v", username, err)
		return "", errors.New("invalid username or password")
	}

	validPassword, err := user.Password.Matches(plaintextPassword, user.Password.Hash)
	if err != nil {
		log.Printf("Error during password comparison for %s: %v", username, err)
		return "", errors.New("error during authentication")
	}
	if !validPassword {
		log.Printf("Login failed for %s: invalid password", username)
		return "", errors.New("invalid username or password")
	}

	log.Printf("User %s authenticated successfully. Generating token.", username)
	token, err := s.generateJWT(user.ID)
	if err != nil {
		log.Printf("Error generating JWT for user %s: %v", username, err)
		return "", errors.New("error generating token")
	}

	return token, nil
}

// generateJWT creates a new JWT token for a given user ID.
func (s *AuthService) generateJWT(userID int) (string, error) {
	claims := jwt.MapClaims{
		"sub": userID,                          // Subject (who the token is for)
		"iat": time.Now().Unix(),               // Issued At
		"exp": time.Now().Add(24*time.Hour).Unix(), // Expiration Time (1 hour)
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	tokenString, err := token.SignedString(jwtSecret)
	if err != nil {
		return "", err
	}

	return tokenString, nil
}
