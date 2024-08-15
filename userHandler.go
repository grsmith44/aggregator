package main

import (
	"encoding/json"
	"net/http"
	"strings"
	"time"

	"github.com/grsmith44/aggregator.git/internal/database"

	"github.com/google/uuid"
)

func (cfg *apiConfig) createUserHandler(w http.ResponseWriter, r *http.Request) {
	type parameters struct {
		Name string `json:"name"`
	}
	decoder := json.NewDecoder(r.Body)
	params := parameters{}
	err := decoder.Decode(&params)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't decode parameters")
		return
	}

	user, err := cfg.DB.CreateUser(r.Context(), database.CreateUserParams{
		ID:        uuid.New(),
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		UserName:  params.Name,
	})
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't create user")
	}

	respondWithJSON(w, http.StatusOK, databaseUserToUser(user))
}

func (cfg *apiConfig) getUserAPIHandler(w http.ResponseWriter, r *http.Request) {
	authHeader := r.Header.Get("Authorization")
	if authHeader == "" {
		respondWithError(w, http.StatusUnauthorized, "Missing Authorization header")
	}

	parts := strings.Split(authHeader, " ")
	if len(parts) != 2 || parts[0] != "ApiKey" {
		respondWithError(w, http.StatusUnauthorized, "Invalid Authorization header format")
	}

	apiKey := parts[1]
	user, err := cfg.DB.GetUserAPI(r.Context(), apiKey)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't find API Key")
	}
	respondWithJSON(w, http.StatusOK, databaseUserToUser(user))
}
