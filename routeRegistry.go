package main

import (
	"fmt"
	"html/template"
	"net/http"
)

func registerRoutes(mux *http.ServeMux, cfg *apiConfig, TESTcfg *apiConfig, tmpl *template.Template) {
	type routes struct {
		method  string
		path    string
		handler http.HandlerFunc
	}
	TESTroutes := []routes{
		{"POST", "/test/users", TESTcfg.createUserHandler},
		{"GET", "/test/users", TESTcfg.middlewareAuth(TESTcfg.getUserAPIHandler)},
		{"POST", "/test/feeds", TESTcfg.middlewareAuth(TESTcfg.createFeedHandler)},
		{"GET", "/test/feeds", TESTcfg.getAllFeedsHandler},
		{"POST", "/test/feed_follows", TESTcfg.middlewareAuth(TESTcfg.createFeedFollowHandler)},
		{"GET", "/test/feed_follows", TESTcfg.middlewareAuth(TESTcfg.getAllFeedFollowsForUserHandler)},
		{"DELETE", "/test/feed_follows/{feedFollowID}", TESTcfg.deleteFeedFollow},
		{"GET", "/test/fetch_full_feed", TESTcfg.fetchRSSFeedHandler},
		{"GET", "/test/posts", TESTcfg.middlewareAuth(TESTcfg.getPostsByUserHandler)},
		{"GET", "/test/healthz", readinessHandler},
		{"GET", "/test/err", errorHandler},
	}
	PRODroutes := []routes{
		{"POST", "/prod/users", cfg.createUserHandler},
		{"GET", "/prod/users", cfg.middlewareAuth(cfg.getUserAPIHandler)},
		{"POST", "/prod/feeds", cfg.middlewareAuth(cfg.createFeedHandler)},
		{"GET", "/prod/feeds", cfg.getAllFeedsHandler},
		{"POST", "/prod/feed_follows", cfg.middlewareAuth(cfg.createFeedFollowHandler)},
		{"GET", "/prod/feed_follows", cfg.middlewareAuth(cfg.getAllFeedFollowsForUserHandler)},
		{"DELETE", "/prod/feed_follows/{feedFollowID}", cfg.deleteFeedFollow},
		{"GET", "/prod/fetch_full_feed", cfg.fetchRSSFeedHandler},
		{"GET", "/prod/posts", cfg.middlewareAuth(cfg.getPostsByUserHandler)},
		{"GET", "/prod/healthz", readinessHandler},
		{"GET", "/prod/err", errorHandler},
	}

	for _, route := range TESTroutes {
		mux.HandleFunc(fmt.Sprintf("%s %s", route.method, route.path), route.handler)
	}

	for _, route := range PRODroutes {
		mux.HandleFunc(fmt.Sprintf("%s %s", route.method, route.path), route.handler)
	}

	mux.HandleFunc("GET /", func(w http.ResponseWriter, r *http.Request) {
		tmpl.ExecuteTemplate(w, "layout.html", nil)
	})

	mux.HandleFunc("POST /login", cfg.loginHandler)
	mux.HandleFunc("GET /dashboard", cfg.middlewareAuth(cfg.dashboardHandler(tmpl)))

	// Special cases
	go mux.HandleFunc("GET /test/start_feed_worker", cfg.startFeedWorker)
	mux.HandleFunc("GET /test/stop_feed_worker", cfg.stopFeedWorker)
	go mux.HandleFunc("GET /prod/start_feed_worker", cfg.startFeedWorker)
	mux.HandleFunc("GET /prod/stop_feed_worker", cfg.stopFeedWorker)
}
