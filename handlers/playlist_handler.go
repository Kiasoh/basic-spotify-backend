package handlers

import (
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	"github.com/kiasoh/basic-spotify-backend/middleware"
	"github.com/kiasoh/basic-spotify-backend/services"
)

type PlaylistHandler struct {
	Service *services.PlaylistService
}

func NewPlaylistHandler(service *services.PlaylistService) *PlaylistHandler {
	return &PlaylistHandler{Service: service}
}

type createPlaylistRequest struct {
	Name string `json:"name"`
}

// Helper to get userID from context
func getUserIDFromContext(r *http.Request) (int, error) {
	userID, ok := r.Context().Value(middleware.UserIDKey).(int)
	if !ok {
		return 0, errors.New("could not retrieve user ID from context")
	}
	return userID, nil
}

func (h *PlaylistHandler) CreatePlaylist(w http.ResponseWriter, r *http.Request) {
	userID, err := getUserIDFromContext(r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusUnauthorized)
		return
	}

	var req createPlaylistRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}
	defer r.Body.Close()

	log.Printf("Handler: User %d creating playlist '%s'", userID, req.Name)
	playlist, err := h.Service.CreatePlaylist(r.Context(), req.Name, userID)
	if err != nil {
		http.Error(w, "Failed to create playlist", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(playlist)
}

func (h *PlaylistHandler) ListUserPlaylists(w http.ResponseWriter, r *http.Request) {
	userID, err := getUserIDFromContext(r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusUnauthorized)
		return
	}

	log.Printf("Handler: User %d listing playlists", userID)
	playlists, err := h.Service.ListUserPlaylists(r.Context(), userID)
	if err != nil {
		http.Error(w, "Failed to list playlists", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(playlists)
}

func (h *PlaylistHandler) AddSongToPlaylist(w http.ResponseWriter, r *http.Request) {
	userID, err := getUserIDFromContext(r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusUnauthorized)
		return
	}

	playlistID, _ := strconv.Atoi(chi.URLParam(r, "playlistID"))
	songID, _ := strconv.Atoi(chi.URLParam(r, "songID"))
	if playlistID == 0 || songID == 0 {
		http.Error(w, "Invalid playlist or song ID", http.StatusBadRequest)
		return
	}

	log.Printf("Handler: User %d adding song %d to playlist %d", userID, songID, playlistID)
	err = h.Service.AddSongToPlaylist(r.Context(), userID, playlistID, songID)
	if err != nil {
		if err.Error() == "forbidden: you do not own this playlist" {
			http.Error(w, err.Error(), http.StatusForbidden)
			return
		}
		http.Error(w, "Failed to add song to playlist", http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

func (h *PlaylistHandler) RemoveSongFromPlaylist(w http.ResponseWriter, r *http.Request) {
	userID, err := getUserIDFromContext(r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusUnauthorized)
		return
	}

	playlistID, _ := strconv.Atoi(chi.URLParam(r, "playlistID"))
	songID, _ := strconv.Atoi(chi.URLParam(r, "songID"))
	if playlistID == 0 || songID == 0 {
		http.Error(w, "Invalid playlist or song ID", http.StatusBadRequest)
		return
	}

	log.Printf("Handler: User %d removing song %d from playlist %d", userID, songID, playlistID)
	err = h.Service.RemoveSongFromPlaylist(r.Context(), userID, playlistID, songID)
	if err != nil {
		if err.Error() == "forbidden: you do not own this playlist" {
			http.Error(w, err.Error(), http.StatusForbidden)
			return
		}
		http.Error(w, "Failed to remove song from playlist", http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

func (h *PlaylistHandler) GetSongsInPlaylist(w http.ResponseWriter, r *http.Request) {
	playlistID, _ := strconv.Atoi(chi.URLParam(r, "playlistID"))
	if playlistID == 0 {
		http.Error(w, "Invalid playlist ID", http.StatusBadRequest)
		return
	}

	log.Printf("Handler: Getting songs for playlist %d", playlistID)
	songs, err := h.Service.GetSongsInPlaylist(r.Context(), playlistID)
	if err != nil {
		http.Error(w, "Failed to get songs in playlist", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(songs)
}
