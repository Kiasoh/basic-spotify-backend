package main

import (
	"context"
	"fmt"
	"log"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/cors"
	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/kiasoh/basic-spotify-backend/handlers"
	"github.com/kiasoh/basic-spotify-backend/middleware"
	"github.com/kiasoh/basic-spotify-backend/repository"
	"github.com/kiasoh/basic-spotify-backend/services"
)

func ConnectSQL() *pgxpool.Pool {
	// TODO: Use environment variables for DSN
	dsn := fmt.Sprintf("postgres://%s:%s@%s:%s/%s",
		"niflheim", "niflguard", "postgres_ds", "5432", "ds_db")

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

func InitRoutes(
	userHandler *handlers.UserHandler,
	authHandler *handlers.AuthHandler,
	songHandler *handlers.SongHandler,
	playlistHandler *handlers.PlaylistHandler,
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

	// Public routes for authentication
	mux.Post("/register", userHandler.Register)
	mux.Post("/login", authHandler.Login)

	// Protected routes
	mux.Group(func(r chi.Router) {
		r.Use(middleware.Auth)

		// Song routes
		r.Get("/songs", songHandler.ListSongs)
		r.Post("/songs", songHandler.CreateSong)
		r.Get("/songs/{id}", songHandler.GetSong)

		// Playlist routes
		r.Get("/playlists", playlistHandler.ListUserPlaylists) // Get playlists for the logged-in user
		r.Post("/playlists", playlistHandler.CreatePlaylist)
		r.Get("/playlists/{playlistID}/songs", playlistHandler.GetSongsInPlaylist)
		r.Post("/playlists/{playlistID}/songs/{songID}", playlistHandler.AddSongToPlaylist)
		r.Delete("/playlists/{playlistID}/songs/{songID}", playlistHandler.RemoveSongFromPlaylist)
	})

	return mux
}

func main() {
	// Initialize database connection
	db := ConnectSQL()
	defer db.Close()

	// --- Initialize Layers ---

	// Repositories
	userRepo := repository.NewUserRepository(db)
	songRepo := repository.NewSongRepository(db)
	playlistRepo := repository.NewPlaylistRepository(db)

	// Services
	userService := services.NewUserService(userRepo)
	authService := services.NewAuthService(userRepo)
	songService := services.NewSongService(songRepo)
	playlistService := services.NewPlaylistService(playlistRepo)

	// Handlers
	userHandler := handlers.NewUserHandler(userService)
	authHandler := handlers.NewAuthHandler(authService)
	songHandler := handlers.NewSongHandler(songService)
	playlistHandler := handlers.NewPlaylistHandler(playlistService)

	// Initialize routes
	router := InitRoutes(userHandler, authHandler, songHandler, playlistHandler)

	server := &http.Server{
		Addr:    ":8081",
		Handler: router,
	}

	log.Println("Server starting on port 8081...")
	if err := server.ListenAndServe(); err != nil {
		log.Fatalf("Server failed to start: %v", err)
	}
}
