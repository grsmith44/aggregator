package main

import (
	"database/sql"
	"log"
	"net/http"
	"os"

	"github.com/grsmith44/aggregator.git/internal/database"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
)

type apiConfig struct {
	DB *database.Queries
}

func dbSetup(dbURL string) apiConfig {
	db, err := sql.Open("postgres", dbURL)
	if err != nil {
		log.Printf("Failed to open database connection with error: %s", err)
		log.Fatal()
	}
	dbQueries := database.New(db)
	apiCfg := apiConfig{
		DB: dbQueries,
	}
	return apiCfg
}

func main() {
	godotenv.Load(".env")
	port := os.Getenv("PORT")
	const filepathRoot = "."
	if port == "" {
		log.Fatal("PORT environment variable is not set")
	}

	dbURL := os.Getenv("DATABASE_URL")
	if dbURL == "" {
		log.Fatal("DATABASE_URL environment variable is not set")
	}

	apiCfg := dbSetup(dbURL)

	mux := http.NewServeMux()

	mux.HandleFunc("POST /v1/users", apiCfg.createUserHandler)
	mux.HandleFunc("GET /v1/users", apiCfg.middlewareAuth(apiCfg.getUserAPIHandler))

	mux.HandleFunc("POST /v1/feeds", apiCfg.middlewareAuth(apiCfg.createFeedHandler))
	mux.HandleFunc("GET /v1/feeds", apiCfg.getAllFeedsHandler)

	mux.HandleFunc("POST /v1/feed_follows", apiCfg.middlewareAuth(apiCfg.createFeedFollowHandler))
	mux.HandleFunc("GET /v1/feed_follows", apiCfg.middlewareAuth(apiCfg.getAllFeedFollowsForUserHandler))
	mux.HandleFunc("DELETE /v1/feed_follows/{feedFollowID}", apiCfg.deleteFeedFollow)

	mux.HandleFunc("GET /v1/fetch_full_feed", apiCfg.fetchRSSFeedHandler)

	go mux.HandleFunc("GET /v1/start_feed_worker", apiCfg.startFeedWorker)
	go mux.HandleFunc("GET /v1/stop_feed_worker", apiCfg.stopFeedWorker)

	mux.HandleFunc("GET /v1/healthz", readinessHandler)
	mux.HandleFunc("GET /v1/err", errorHandler)

	server := http.Server{
		Addr:    ":" + port,
		Handler: mux,
	}

	log.Printf("Serving on port: %s\n", port)
	log.Fatal(server.ListenAndServe())
}
