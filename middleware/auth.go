package middleware

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"strings"

	"github.com/golang-jwt/jwt/v5"
)

// TODO: Move this to a secure location like environment variables
var jwtSecret = []byte("supersecretkey")

type contextKey string

const UserIDKey contextKey = "userID"

func Auth(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" {
			log.Println("Auth error: Authorization header missing")
			http.Error(w, "Authorization header required", http.StatusUnauthorized)
			return
		}

		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || parts[0] != "Bearer" {
			log.Println("Auth error: Invalid Authorization header format")
			http.Error(w, "Invalid Authorization header format", http.StatusUnauthorized)
			return
		}

		tokenString := parts[1]

		token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
			}
			return jwtSecret, nil
		})

		if err != nil {
			log.Printf("Auth error: Invalid token: %v", err)
			http.Error(w, "Invalid token", http.StatusUnauthorized)
			return
		}

		if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
			// Extract user ID from the 'sub' claim
			if userIDFloat, ok := claims["sub"].(float64); ok {
				userID := int(userIDFloat)
				// Add user ID to the request context
				ctx := context.WithValue(r.Context(), UserIDKey, userID)
				log.Printf("Authenticated user with ID: %d", userID)
				next.ServeHTTP(w, r.WithContext(ctx))
				return
			}
		}

		log.Println("Auth error: Invalid claims or token")
		http.Error(w, "Invalid token", http.StatusUnauthorized)
	})
}
