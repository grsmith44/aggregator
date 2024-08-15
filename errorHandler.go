package main

import (
	"net/http"
)

func errorHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	respondWithError(w, 500, "Internal Server Error")
}
