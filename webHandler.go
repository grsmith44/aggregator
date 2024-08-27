package main

import (
	"html/template"
	"net/http"

	"github.com/google/uuid"
	"github.com/grsmith44/aggregator.git/internal/database"
)

func (cfg *apiConfig) loginHandler(w http.ResponseWriter, r *http.Request) {
	username := r.FormValue("username")
	userAPIKey, err := cfg.DB.GetUserAPIKeyByName(r.Context(), username)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Error getting user API from user name")
	}
	r.Header.Add("Authorization", userAPIKey)
	http.Redirect(w, r, "/dashboard", http.StatusSeeOther)
}

func (cfg *apiConfig) dashboardHandler(tmpl *template.Template) authedHandler {
	return func(w http.ResponseWriter, r *http.Request, user database.User) {
		username, _ := r.Cookie("username")
		userAPIString, err := cfg.DB.GetUserAPIKeyByName(r.Context(), username.Value)
		if err != nil {
			respondWithError(w, http.StatusInternalServerError, "Error getting user API from user name")
		}
		userAPI, err := uuid.Parse(userAPIString)
		if err != nil {
			respondWithError(w, http.StatusInternalServerError, "Error parsing UUID for provided username")
			return
		}

		DBavailableFeeds, err := cfg.DB.SelectAllFeeds(r.Context())
		if err != nil {
			respondWithError(w, http.StatusInternalServerError, "Error fetching feeds")
			return
		}
		availableFeeds := batchDatabaseFeedToFeeds(DBavailableFeeds)

		followedFeeds, err := cfg.DB.GetAllFeedNamesFollowedByUser(r.Context(), userAPI)
		if err != nil {
			respondWithError(w, http.StatusInternalServerError, "Error fetching followed feeds")
			return
		}

		DBrecentPosts, err := cfg.DB.GetPostsByUser(r.Context(), database.GetPostsByUserParams{})
		if err != nil {
			respondWithError(w, http.StatusInternalServerError, "Error fetching recent posts")
			return
		}
		recentPosts := batchPostToPost(DBrecentPosts)

		data := struct {
			Username       string
			AvailableFeeds []Feed
			FollowedFeeds  []string
			RecentPosts    []Post
		}{
			Username:       username.Value,
			AvailableFeeds: availableFeeds,
			FollowedFeeds:  followedFeeds,
			RecentPosts:    recentPosts,
		}
		tmpl.ExecuteTemplate(w, "layout.html", data)
	}
}
