package main

import (
	"database/sql"
	"html/template"
	"log"
	"net/http"
	"os"

	"github.com/grsmith44/aggregator.git/internal/database"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
	"github.com/pressly/goose/v3"
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
	if err := goose.Up(db, "sql/schema"); err != nil {
		log.Fatal(err)
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
	TESTdbURL := os.Getenv("TEST_DATABASE_URL")
	if TESTdbURL == "" {
		log.Fatal("TEST_DATABASE_URL environment variable is not set")
	}
	apiCfg := dbSetup(dbURL)
	apiCfgTEST := dbSetup(TESTdbURL)

	tmpl := template.Must(template.ParseGlob("templates/*.html"))

	mux := http.NewServeMux()

	registerRoutes(mux, &apiCfg, &apiCfgTEST, tmpl)

	server := http.Server{
		Addr:    ":" + port,
		Handler: mux,
	}

	log.Printf("Serving on port: %s\n", port)
	log.Fatal(server.ListenAndServe())
}
