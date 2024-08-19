package main

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/grsmith44/aggregator.git/internal/database"

	"github.com/google/uuid"
)

func (cfg *apiConfig) createFeedFollowHandler(w http.ResponseWriter, r *http.Request, user database.User) {
	type parameters struct {
		FeedId uuid.UUID `json:"feed_id"`
	}
	decoder := json.NewDecoder(r.Body)
	params := parameters{}
	err := decoder.Decode(&params)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't decode parameters")
		return
	}

	feedFollow, err := cfg.DB.CreateFeedFollow(r.Context(), database.CreateFeedFollowParams{
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		ID:        uuid.New(),
		FeedID:    params.FeedId,
		UserID:    user.ID,
	})
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't create feed")
	}

	respondWithJSON(w, http.StatusOK, databaseFeedFollowToFeedFollow(feedFollow))
}

func (cfg *apiConfig) deleteFeedFollow(w http.ResponseWriter, r *http.Request) {
	uuid, err := uuid.Parse(r.PathValue("feedFollowID"))
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid URL for feed follow deletion")
	}

	result, err := cfg.DB.DeleteFeedFollow(r.Context(), uuid)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Failed to delete feed follow")
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Error checking affected row")
	}
	if rowsAffected == 0 {
		respondWithError(w, http.StatusNotFound, "Feed follow not found")
	}

	w.WriteHeader(http.StatusOK)
}

func (cfg *apiConfig) getAllFeedFollowsForUserHandler(w http.ResponseWriter, r *http.Request, user database.User) {
	dbFeedFollows, err := cfg.DB.GetAllFeedFollowsForUser(r.Context(), user.ID)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Failed to retrieve all feeds")
	}
	feedFollows := batchFeedFollowsToFeedFollows(dbFeedFollows)
	respondWithJSON(w, http.StatusOK, feedFollows)
}
