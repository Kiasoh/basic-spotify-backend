package handlers

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/kiasoh/basic-spotify-backend/services"
)

type UserHandler struct {
	Service *services.UserService
}

func NewUserHandler(service *services.UserService) *UserHandler {
	return &UserHandler{Service: service}
}

type registerRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

func (h *UserHandler) Register(w http.ResponseWriter, r *http.Request) {
	var req registerRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		log.Printf("Error decoding registration request: %v", err)
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}
	defer r.Body.Close()

	log.Printf("Handling registration request for user: %s", req.Username)

	user, err := h.Service.RegisterUser(r.Context(), req.Username, req.Password)
	if err != nil {
		// This is a simple error mapping. In a real app, you might use custom error types.
		if err.Error() == "user with this username already exists" {
			http.Error(w, err.Error(), http.StatusConflict)
			return
		}
		if err.Error() == "password must be at least 8 characters long" {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		http.Error(w, "Failed to register user", http.StatusInternalServerError)
		return
	}

	// Exclude password from the response
	user.Password.Hash = nil
	user.Password.Plaintext = nil

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	if err := json.NewEncoder(w).Encode(user); err != nil {
		log.Printf("Error encoding registration response: %v", err)
	}
}
