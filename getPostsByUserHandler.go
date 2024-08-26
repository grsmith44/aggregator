package main

import (
	"net/http"

	"github.com/grsmith44/aggregator.git/internal/database"
)

func (cfg *apiConfig) getPostsByUserHandler(w http.ResponseWriter, r *http.Request, user database.User) {
	posts, err := cfg.DB.GetPostsByUser(r.Context(), user.ID)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't retrieve feeds for that user")
	}

	postList := batchPostToPost(posts)
	respondWithJSON(w, http.StatusOK, postList)
}
