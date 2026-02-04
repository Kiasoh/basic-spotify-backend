package handlers

import (
	"encoding/json"
	"log"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	"github.com/kiasoh/basic-spotify-backend/services"
)

type PlaylistHandler struct {
	Service *services.PlaylistService
}

func NewPlaylistHandler(service *services.PlaylistService) *PlaylistHandler {
	return &PlaylistHandler{Service: service}
}

type createPlaylistRequest struct {
	Name        string  `json:"name"`
	Description *string `json:"description"`
}

type updatePlaylistRequest struct {
	Name        string  `json:"name"`
	Description *string `json:"description"`
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
	playlist, err := h.Service.CreatePlaylist(r.Context(), userID, req.Name, req.Description)
	if err != nil {
		http.Error(w, "Failed to create playlist", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(playlist)
}

func (h *PlaylistHandler) UpdatePlaylistDetails(w http.ResponseWriter, r *http.Request) {
	userID, err := getUserIDFromContext(r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusUnauthorized)
		return
	}

	playlistID, _ := strconv.Atoi(chi.URLParam(r, "playlistID"))
	if playlistID == 0 {
		http.Error(w, "Invalid playlist ID", http.StatusBadRequest)
		return
	}

	var req updatePlaylistRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}
	defer r.Body.Close()

	log.Printf("Handler: User %d updating playlist %d", userID, playlistID)
	playlist, err := h.Service.UpdatePlaylistDetails(r.Context(), userID, playlistID, req.Name, req.Description)
	if err != nil {
		if err.Error() == "forbidden: you do not own this playlist" || err.Error() == "forbidden: this playlist is not modifiable" {
			http.Error(w, err.Error(), http.StatusForbidden)
			return
		}
		http.Error(w, "Failed to update playlist", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
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

func (h *PlaylistHandler) AddTrackToPlaylist(w http.ResponseWriter, r *http.Request) {
	userID, err := getUserIDFromContext(r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusUnauthorized)
		return
	}

	playlistID, _ := strconv.Atoi(chi.URLParam(r, "playlistID"))
	trackID := chi.URLParam(r, "trackID")
	if playlistID == 0 || trackID == "" {
		http.Error(w, "Invalid playlist or track ID", http.StatusBadRequest)
		return
	}

	log.Printf("Handler: User %d adding track %s to playlist %d", userID, trackID, playlistID)
	err = h.Service.AddTrackToPlaylist(r.Context(), userID, playlistID, trackID)
	if err != nil {
		if err.Error() == "forbidden: you do not own this playlist" || err.Error() == "forbidden: this playlist is not modifiable" {
			http.Error(w, err.Error(), http.StatusForbidden)
			return
		}
		http.Error(w, "Failed to add track to playlist", http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

func (h *PlaylistHandler) RemoveTrackFromPlaylist(w http.ResponseWriter, r *http.Request) {
	userID, err := getUserIDFromContext(r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusUnauthorized)
		return
	}

	playlistID, _ := strconv.Atoi(chi.URLParam(r, "playlistID"))
	trackID := chi.URLParam(r, "trackID")
	if playlistID == 0 || trackID == "" {
		http.Error(w, "Invalid playlist or track ID", http.StatusBadRequest)
		return
	}

	log.Printf("Handler: User %d removing track %s from playlist %d", userID, trackID, playlistID)
	err = h.Service.RemoveTrackFromPlaylist(r.Context(), userID, playlistID, trackID)
	if err != nil {
		if err.Error() == "forbidden: you do not own this playlist" || err.Error() == "forbidden: this playlist is not modifiable" {
			http.Error(w, err.Error(), http.StatusForbidden)
			return
		}
		http.Error(w, "Failed to remove track from playlist", http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

func (h *PlaylistHandler) GetTracksInPlaylist(w http.ResponseWriter, r *http.Request) {
	playlistID, _ := strconv.Atoi(chi.URLParam(r, "playlistID"))
	if playlistID == 0 {
		http.Error(w, "Invalid playlist ID", http.StatusBadRequest)
		return
	}

	log.Printf("Handler: Getting tracks for playlist %d", playlistID)
	tracks, err := h.Service.GetTracksInPlaylist(r.Context(), playlistID)
	if err != nil {
		http.Error(w, "Failed to get tracks in playlist", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(tracks)
}
