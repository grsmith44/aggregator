package main

import (
	"net/http"
	"strings"

	"github.com/grsmith44/aggregator.git/internal/auth"
	"github.com/grsmith44/aggregator.git/internal/database"
)

type authedHandler func(http.ResponseWriter, *http.Request, database.User)

func (cfg *apiConfig) middlewareAuth(handler authedHandler) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		apiKey, err := auth.GetApiKey(r.Header)
		if err != nil {
			respondWithError(w, http.StatusUnauthorized, "Missing Authorization header or ApiKey")
		}

		apiKey = strings.TrimPrefix(apiKey, "ApiKey ")
		user, err := cfg.DB.GetUserByAPI(r.Context(), apiKey)
		if err != nil {
			respondWithError(w, http.StatusInternalServerError, "Couldn't find user")
		}
		handler(w, r, user)
	}
}
