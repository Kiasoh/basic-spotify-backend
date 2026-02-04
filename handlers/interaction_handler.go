package handlers

import (
	"encoding/json"
	"errors"
	"log"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/kiasoh/basic-spotify-backend/middleware"
	"github.com/kiasoh/basic-spotify-backend/services"
)

type InteractionHandler struct {
	Service *services.InteractionService
}

func NewInteractionHandler(service *services.InteractionService) *InteractionHandler {
	return &InteractionHandler{Service: service}
}

type createInteractionRequest struct {
	Type string `json:"type"` // e.g., "like", "play"
}

// NOTE: This helper is duplicated from playlist_handler.go. It could be moved to a shared package.
func getUserIDFromContext(r *http.Request) (int, error) {
	userID, ok := r.Context().Value(middleware.UserIDKey).(int)
	if !ok {
		return 0, errors.New("could not retrieve user ID from context")
	}
	return userID, nil
}

func (h *InteractionHandler) CreateInteraction(w http.ResponseWriter, r *http.Request) {
	userID, err := getUserIDFromContext(r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusUnauthorized)
		return
	}

	trackID := chi.URLParam(r, "trackID")
	if trackID == "" {
		http.Error(w, "Invalid track ID", http.StatusBadRequest)
		return
	}

	var req createInteractionRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}
	defer r.Body.Close()

	log.Printf("Handler: User %d creating interaction '%s' for track %s", userID, req.Type, trackID)

	err = h.Service.CreateInteraction(r.Context(), userID, trackID, req.Type)
	if err != nil {
		http.Error(w, "Failed to create interaction", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
}

func (h *InteractionHandler) GetInteractionsForTrack(w http.ResponseWriter, r *http.Request) {
	trackID := chi.URLParam(r, "trackID")
	if trackID == "" {
		http.Error(w, "Invalid track ID", http.StatusBadRequest)
		return
	}

	log.Printf("Handler: Getting interactions for track %s", trackID)
	interactions, err := h.Service.GetInteractionsForTrack(r.Context(), trackID)
	if err != nil {
		http.Error(w, "Failed to get interactions", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(interactions)
}
