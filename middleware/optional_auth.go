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

// OptionalAuth attempts to authenticate a user from a JWT token.
// If successful, the userID is added to the request context.
// If no token is provided, or the token is invalid, it proceeds to the next handler
// without setting the userID and without returning an error (i.e., it doesn't block unauthenticated requests).
func OptionalAuth(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" {
			// No Authorization header, proceed without user ID
			next.ServeHTTP(w, r)
			return
		}

		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || parts[0] != "Bearer" {
			// Invalid header format, proceed without user ID
			log.Println("OptionalAuth: Invalid Authorization header format, proceeding unauthenticated")
			next.ServeHTTP(w, r)
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
			// Invalid token, proceed without user ID
			log.Printf("OptionalAuth: Invalid token: %v, proceeding unauthenticated", err)
			next.ServeHTTP(w, r)
			return
		}

		if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
			// Extract user ID from the 'sub' claim
			if userIDFloat, ok := claims["sub"].(float64); ok {
				userID := int(userIDFloat)
				// Add user ID to the request context
				ctx := context.WithValue(r.Context(), UserIDKey, userID)
				log.Printf("OptionalAuth: Authenticated user with ID: %d", userID)
				next.ServeHTTP(w, r.WithContext(ctx))
				return
			}
		}

		// Invalid claims or token, proceed without user ID
		log.Println("OptionalAuth: Invalid claims or token, proceeding unauthenticated")
		next.ServeHTTP(w, r)
	})
}
