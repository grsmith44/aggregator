package main

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/grsmith44/aggregator/internal/database"

	"github.com/google/uuid"
)

func (cfg *configApi) createUserHandler(w http.ResponseWriter, r *http.Request) {
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

	user, err := cfg.DB.CreateUser(r.context(), database.CreateUserParams{
		ID:         uuid.New(),
		CREATED_AT: time.Now(),
		UPDATED_AT: time.Now(),
		USER_NAME:  params.Name,
	})
}
