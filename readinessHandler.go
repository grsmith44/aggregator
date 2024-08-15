package main

import (
	"net/http"
)

func readinessHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	type readinessStruct struct {
		Status string `json:"status"`
	}
	respondWithJSON(w, http.StatusOK, readinessStruct{
		Status: "ok",
	})
}
