package handlers

import (
	"encoding/json"
	"log"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
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

func (h *InteractionHandler) CreateInteraction(w http.ResponseWriter, r *http.Request) {
	userID, err := getUserIDFromContext(r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusUnauthorized)
		return
	}

	songID, err := strconv.Atoi(chi.URLParam(r, "songID"))
	if err != nil {
		http.Error(w, "Invalid song ID", http.StatusBadRequest)
		return
	}

	var req createInteractionRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}
	defer r.Body.Close()

	log.Printf("Handler: User %d creating interaction '%s' for song %d", userID, req.Type, songID)

	err = h.Service.CreateInteraction(r.Context(), userID, songID, req.Type)
	if err != nil {
		http.Error(w, "Failed to create interaction", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
}

func (h *InteractionHandler) GetInteractionsForSong(w http.ResponseWriter, r *http.Request) {
	songID, err := strconv.Atoi(chi.URLParam(r, "songID"))
	if err != nil {
		http.Error(w, "Invalid song ID", http.StatusBadRequest)
		return
	}

	log.Printf("Handler: Getting interactions for song %d", songID)
	interactions, err := h.Service.GetInteractionsForSong(r.Context(), songID)
	if err != nil {
		http.Error(w, "Failed to get interactions", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(interactions)
}
