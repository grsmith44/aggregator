package main

import (
	"encoding/json"
	"net/http"

	"github.com/grsmith44/aggregator.git/internal/database"
)

func (cfg *apiConfig) getPostsByUserHandler(w http.ResponseWriter, r *http.Request, user database.User) {
	type parameters struct {
		Limit *int `json:"limit,omitempty"`
	}
	decoder := json.NewDecoder(r.Body)
	params := parameters{}
	err := decoder.Decode(&params)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't decode parameters ")
		return
	}

	var limit *int
	if params.Limit != nil {
		l := int(*params.Limit)
		limit = &l
	}

	posts, err := cfg.DB.GetPostsByUser(r.Context(), database.GetPostsByUserParams{
		UserID:  user.ID,
		Column2: limit,
	})
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't retrieve feeds for that user")
	}

	postList := batchPostToPost(posts)
	respondWithJSON(w, http.StatusOK, postList)
}
