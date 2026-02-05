package main

import (
	"context"
	"fmt"
	"log"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/cors"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/segmentio/kafka-go"

	"github.com/kiasoh/basic-spotify-backend/handlers"
	"github.com/kiasoh/basic-spotify-backend/middleware"
	"github.com/kiasoh/basic-spotify-backend/repository"
	"github.com/kiasoh/basic-spotify-backend/services"
)

func ConnectSQL() *pgxpool.Pool {
	dsn := fmt.Sprintf("postgres://%s:%s@%s:%s/%s",
		"niflheim", "niflguard", "postgres_ds", "5432", "dsdb")

	poolconfig, err := pgxpool.ParseConfig(dsn)
	if err != nil {
		log.Fatalf("Unable to parse DSN: %v", err)
	}

	pool, err := pgxpool.NewWithConfig(context.Background(), poolconfig)
	if err != nil {
		log.Fatalf("Unable to connect to database: %v", err)
	}

	if err = pool.Ping(context.Background()); err != nil {
		log.Fatalf("Unable to ping database: %v", err)
	}
	log.Println("Database connection established")
	return pool
}

func InitKafka() *kafka.Writer {
	// TODO: Use environment variables for Kafka URL
	kafkaURL := "194.147.142.26:9094"
	writer := &kafka.Writer{
		Addr:     kafka.TCP(kafkaURL),
		Topic:    "interactions",
		Balancer: &kafka.LeastBytes{},
	}
	log.Println("Kafka writer initialized for topic 'interactions'")
	return writer
}

func InitRoutes(
	userHandler *handlers.UserHandler,
	authHandler *handlers.AuthHandler,
	trackHandler *handlers.SpotifyTrackHandler,
	playlistHandler *handlers.PlaylistHandler,
	interactionHandler *handlers.InteractionHandler,
) http.Handler {
	mux := chi.NewRouter()

	mux.Use(cors.Handler(cors.Options{
		AllowedOrigins:   []string{"*"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS", "PATCH"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type"},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: true,
		MaxAge:           300,
	}))

	// Public routes
	mux.Post("/register", userHandler.Register)
	mux.Post("/login", authHandler.Login)
	mux.Get("/tracks/{trackID}", trackHandler.GetByTrackID)
	mux.Get("/tracks", trackHandler.ListTracks)
	mux.Get("/tracks/search", trackHandler.SearchTracks) // New Search Route
	// Protected routes
	mux.Group(func(r chi.Router) {
		r.Use(middleware.Auth)

		// Interaction routes
		r.Post("/tracks/{trackID}/interact", interactionHandler.CreateInteraction)
		r.Get("/tracks/{trackID}/interactions", interactionHandler.GetInteractionsForTrack)

		// Playlist routes
		r.Get("/playlists", playlistHandler.ListUserPlaylists)
		r.Post("/playlists", playlistHandler.CreatePlaylist)
		r.Put("/playlists/{playlistID}", playlistHandler.UpdatePlaylistDetails)
		r.Get("/playlists/{playlistID}/tracks", playlistHandler.GetTracksInPlaylist)
		r.Post("/playlists/{playlistID}/tracks/{trackID}", playlistHandler.AddTrackToPlaylist)
		r.Delete("/playlists/{playlistID}/tracks/{trackID}", playlistHandler.RemoveTrackFromPlaylist)
	})

	return mux
}

func main() {
	// Initialize database connection
	db := ConnectSQL()
	defer db.Close()

	// Initialize Kafka Writer
	kafkaWriter := InitKafka()
	defer kafkaWriter.Close()

	// --- Initialize Layers ---

	// Repositories
	userRepo := repository.NewUserRepository(db)
	trackRepo := repository.NewSpotifyTrackRepository(db)
	playlistRepo := repository.NewPlaylistRepository(db)
	interactionRepo := repository.NewInteractionRepository(db)

	// Services
	interactionService := services.NewInteractionService(interactionRepo, kafkaWriter, trackRepo, userRepo)
	userService := services.NewUserService(db, userRepo, playlistRepo)
	authService := services.NewAuthService(userRepo)
	trackService := services.NewSpotifyTrackService(trackRepo, interactionService)
	playlistService := services.NewPlaylistService(playlistRepo, interactionService)

	// Handlers
	userHandler := handlers.NewUserHandler(userService, authService)
	authHandler := handlers.NewAuthHandler(authService)
	trackHandler := handlers.NewSpotifyTrackHandler(trackService)
	playlistHandler := handlers.NewPlaylistHandler(playlistService)
	interactionHandler := handlers.NewInteractionHandler(interactionService)

	// Initialize routes
	router := InitRoutes(userHandler, authHandler, trackHandler, playlistHandler, interactionHandler)

	server := &http.Server{
		Addr:    ":8081",
		Handler: router,
	}

	log.Println("Server starting on port 8081...")
	if err := server.ListenAndServe(); err != nil {
		log.Fatalf("Server failed to start: %v", err)
	}
}
