package handlers

import (
	"encoding/json"
	"log"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	"github.com/kiasoh/basic-spotify-backend/services"
)

type SpotifyTrackHandler struct {
	Service *services.SpotifyTrackService
}

func NewSpotifyTrackHandler(service *services.SpotifyTrackService) *SpotifyTrackHandler {
	return &SpotifyTrackHandler{Service: service}
}

func (h *SpotifyTrackHandler) GetByTrackID(w http.ResponseWriter, r *http.Request) {
	trackID := chi.URLParam(r, "trackID")

	log.Printf("Handler: Handling get track request for ID: %s", trackID)

	track, err := h.Service.GetByTrackID(r.Context(), trackID)
	if err != nil {
		http.Error(w, "Track not found", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(track)
}

func (h *SpotifyTrackHandler) ListTracks(w http.ResponseWriter, r *http.Request) {
	limitStr := r.URL.Query().Get("limit")
	offsetStr := r.URL.Query().Get("offset")

	limit, err := strconv.Atoi(limitStr)
	if err != nil || limit <= 0 {
		limit = 20 // Default limit
	}

	offset, err := strconv.Atoi(offsetStr)
	if err != nil || offset < 0 {
		offset = 0 // Default offset
	}

	log.Printf("Handler: Handling list tracks request with limit %d and offset %d", limit, offset)

	tracks, err := h.Service.List(r.Context(), limit, offset)
	if err != nil {
		http.Error(w, "Failed to list tracks", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(tracks)
}
