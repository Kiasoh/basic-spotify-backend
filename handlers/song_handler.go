package handlers

import (
	"encoding/json"
	"log"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	"github.com/kiasoh/basic-spotify-backend/services"
)

type SongHandler struct {
	Service *services.SongService
}

func NewSongHandler(service *services.SongService) *SongHandler {
	return &SongHandler{Service: service}
}

type createSongRequest struct {
	Name   string `json:"name"`
	Artist string `json:"artist"`
	Album  string `json:"album"`
	Genre  string `json:"genre"`
}

func (h *SongHandler) CreateSong(w http.ResponseWriter, r *http.Request) {
	var req createSongRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		log.Printf("Handler: Error decoding create song request: %v", err)
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}
	defer r.Body.Close()

	log.Printf("Handler: Handling create song request for: %s", req.Name)

	song, err := h.Service.CreateSong(r.Context(), req.Name, req.Artist, req.Album, req.Genre)
	if err != nil {
		http.Error(w, "Failed to create song", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(song)
}

func (h *SongHandler) GetSong(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "Invalid song ID", http.StatusBadRequest)
		return
	}

	log.Printf("Handler: Handling get song request for ID: %d", id)

	song, err := h.Service.GetSong(r.Context(), id)
	if err != nil {
		// In a real app, you'd check for a specific "not found" error
		http.Error(w, "Song not found", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(song)
}

func (h *SongHandler) ListSongs(w http.ResponseWriter, r *http.Request) {
	log.Println("Handler: Handling list songs request")

	songs, err := h.Service.ListSongs(r.Context())
	if err != nil {
		http.Error(w, "Failed to list songs", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(songs)
}
