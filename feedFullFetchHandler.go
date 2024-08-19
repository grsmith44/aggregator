package main

import (
	"encoding/json"
	"encoding/xml"
	"fmt"
	"net/http"
	"time"

	"github.com/google/uuid"
)

// RSS is the root element of the XML feed
type RSS struct {
	XMLName xml.Name `xml:"rss"`
	Channel Channel  `xml:"channel"`
}

// Channel contains metadata about the feed and a list of items
type Channel struct {
	Title       string `xml:"title"`
	Link        string `xml:"link"`
	Description string `xml:"description"`
	Items       []Item `xml:"item"`
}

// Item represents a single article in the feed
type Item struct {
	Title       string `xml:"title"`
	Link        string `xml:"link"`
	Description string `xml:"description"`
	PubDate     string `xml:"pubDate"`
}

func (cfg *apiConfig) fetchRSSFeedHandler(w http.ResponseWriter, r *http.Request) {
	type parameters struct {
		FeedId uuid.UUID `json:"feed_id"`
	}
	decoder := json.NewDecoder(r.Body)
	params := parameters{}
	err := decoder.Decode(&params)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't decode parameters to fetch RSS")
	}

	dbFeed, err := cfg.DB.GetFeedFromID(r.Context(), params.FeedId)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldnt find database entry with provided feedID")
	}

	feed := databaseFeedToFeed(dbFeed)

	client := &http.Client{
		Timeout: 10 * time.Second,
	}

	resp, err := client.Get(feed.Url)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, fmt.Sprintf("Error fetching the RSS: %s", err))
	}
	defer resp.Body.Close()

	rss := RSS{}
	xmlDecoder := xml.NewDecoder(resp.Body)
	err = xmlDecoder.Decode(&rss)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, fmt.Sprintf("Error decoding the RSS: %s", err))
	}
	cfg.DB.MarkFeedAsFetched(r.Context(), feed.ID)
	respondWithJSON(w, http.StatusOK, rss)
}
