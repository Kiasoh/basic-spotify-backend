package handlers

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	"github.com/kiasoh/basic-spotify-backend/middleware" // Added
	"github.com/kiasoh/basic-spotify-backend/models"     // Added
	"github.com/kiasoh/basic-spotify-backend/services"
)

type SpotifyTrackHandler struct {
	Service *services.SpotifyTrackService
}

func NewSpotifyTrackHandler(service *services.SpotifyTrackService) *SpotifyTrackHandler {
	return &SpotifyTrackHandler{Service: service}
}

// enrichTracksWithInteractionState processes a slice of SpotifyTrack and returns SpotifyTrackResponse
// with interaction states if a user ID is present in the request context.
func (h *SpotifyTrackHandler) enrichTracksWithInteractionState(r *http.Request, tracks []models.SpotifyTrack) ([]models.SpotifyTrackResponse, error) {
	trackResponses := make([]models.SpotifyTrackResponse, len(tracks))
	userID, _ := r.Context().Value(middleware.UserIDKey).(int) // Get userID, 0 if not present

	if userID != 0 { // User is authenticated, fetch interaction states
		trackIDs := make([]string, len(tracks))
		for i, track := range tracks {
			trackIDs[i] = track.TrackID
		}

		interactionStates, err := h.Service.InteractionService.GetTrackInteractionStates(r.Context(), userID, trackIDs)
		if err != nil {
			log.Printf("Handler: Error getting interaction states for user %d: %v", userID, err)
			// Continue without interaction states if there's an error
		}

		for i, track := range tracks {
			trackResponses[i].SpotifyTrack = track
			if state, ok := interactionStates[track.TrackID]; ok {
				trackResponses[i].InteractionState = state
			} else {
				trackResponses[i].InteractionState = models.TrackStateNeutral // Default to neutral
			}
		}
	} else { // User is not authenticated, just populate SpotifyTrack
		for i, track := range tracks {
			trackResponses[i].SpotifyTrack = track
			// InteractionState will be omitted due to omitempty tag
		}
	}
	return trackResponses, nil
}

func (h *SpotifyTrackHandler) GetByTrackID(w http.ResponseWriter, r *http.Request) {
	trackID := chi.URLParam(r, "trackID")

	log.Printf("Handler: Handling get track request for ID: %s", trackID)

	track, err := h.Service.GetByTrackID(r.Context(), trackID)
	if err != nil {
		http.Error(w, "Track not found", http.StatusNotFound)
		return
	}

	// Prepare response with interaction states if user is authenticated
	var trackResponse models.SpotifyTrackResponse
	trackResponse.SpotifyTrack = *track // Dereference the pointer

	userID, _ := r.Context().Value(middleware.UserIDKey).(int) // Get userID, 0 if not present

	if userID != 0 { // User is authenticated, fetch interaction state
		interactionStates, err := h.Service.InteractionService.GetTrackInteractionStates(r.Context(), userID, []string{trackID})
		if err != nil {
			log.Printf("Handler: Error getting interaction state for user %d and track %s: %v", userID, trackID, err)
			// Continue without interaction state if there's an error
		}

		if state, ok := interactionStates[trackID]; ok {
			trackResponse.InteractionState = state
		} else {
			trackResponse.InteractionState = models.TrackStateNeutral // Default to neutral
		}
	}
	// If user is not authenticated, InteractionState will be omitted due to omitempty tag

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(trackResponse)
}

func (h *SpotifyTrackHandler) ListTracks(w http.ResponseWriter, r *http.Request) {
	limitStr := r.URL.Query().Get("limit")
	offsetStr := r.URL.Query().Get("offset")
	sortBy := r.URL.Query().Get("sort_by") // New parameter
	order := r.URL.Query().Get("order")

	limit, err := strconv.Atoi(limitStr)
	if err != nil || limit <= 0 {
		limit = 20 // Default limit
	}

	offset, err := strconv.Atoi(offsetStr)
	if err != nil || offset < 0 {
		offset = 0 // Default offset
	}

	log.Printf("Handler: Handling list tracks request with limit %d, offset %d, and sort by %s", limit, offset, sortBy)

	tracks, err := h.Service.List(r.Context(), limit, offset, sortBy, order)
	if err != nil {
		// Handle invalid sort field error
		if err.Error() == fmt.Sprintf("invalid sort field: %s", sortBy) {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		http.Error(w, "Failed to list tracks", http.StatusInternalServerError)
		return
	}

	trackResponses, err := h.enrichTracksWithInteractionState(r, tracks)
	if err != nil {
		log.Printf("Handler: Error enriching tracks with interaction state in ListTracks: %v", err)
		http.Error(w, "Failed to process tracks", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(trackResponses)
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

	trackResponses, err := h.enrichTracksWithInteractionState(r, tracks)
	if err != nil {
		log.Printf("Handler: Error enriching tracks with interaction state in SearchTracks: %v", err)
		http.Error(w, "Failed to process tracks", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(trackResponses)
}
