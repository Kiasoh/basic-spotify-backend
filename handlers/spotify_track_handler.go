package handlers

import (
	"encoding/json"
	"fmt"
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
	sortBy := r.URL.Query().Get("sort_by") // New parameter

	limit, err := strconv.Atoi(limitStr)
	if err != nil || limit <= 0 {
		limit = 20 // Default limit
	}

	offset, err := strconv.Atoi(offsetStr)
	if err != nil || offset < 0 {
		offset = 0 // Default offset
	}

	log.Printf("Handler: Handling list tracks request with limit %d, offset %d, and sort by %s", limit, offset, sortBy)

	tracks, err := h.Service.List(r.Context(), limit, offset, sortBy)
	if err != nil {
		// Handle invalid sort field error
		if err.Error() == fmt.Sprintf("invalid sort field: %s", sortBy) {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		http.Error(w, "Failed to list tracks", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(tracks)
}

func (h *SpotifyTrackHandler) SearchTracks(w http.ResponseWriter, r *http.Request) {
	query := r.URL.Query().Get("q")
	searchField := r.URL.Query().Get("field")
	limitStr := r.URL.Query().Get("limit")
	offsetStr := r.URL.Query().Get("offset")

	if query == "" || searchField == "" {
		http.Error(w, "Query (q) and search field (field) are required", http.StatusBadRequest)
		return
	}

	limit, err := strconv.Atoi(limitStr)
	if err != nil || limit <= 0 {
		limit = 20 // Default limit
	}

	offset, err := strconv.Atoi(offsetStr)
	if err != nil || offset < 0 {
		offset = 0 // Default offset
	}

	log.Printf("Handler: Handling search tracks request for query '%s' in field '%s' with limit %d and offset %d", query, searchField, limit, offset)

	tracks, err := h.Service.Search(r.Context(), query, searchField, limit, offset)
	if err != nil {
		// Specific error for invalid field from repo/service
		if err.Error() == fmt.Sprintf("invalid search field: %s", searchField) {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		http.Error(w, "Failed to search tracks", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(tracks)
}
