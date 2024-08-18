package main

import (
	"context"
	"encoding/json"
	"encoding/xml"
	"log"
	"net/http"
	"sync"
	"time"
)

var (
	workerCtx    context.Context
	workerCancel context.CancelFunc
)

type config struct {
	BatchSize int `json:"batch_size"`
	Interval  int `json:"interval"`
}

func (cfg *apiConfig) startFeedWorker(w http.ResponseWriter, r *http.Request) {
	log.Println("Starting the feed worker")
	params := config{}
	err := json.NewDecoder(r.Body).Decode(&params)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Couldn't decode parameters for interval or batch size")
		return
	}
	log.Printf("Parameters succesfully parsed: batch size %d, interval: %d s \n", params.BatchSize, params.Interval)

	if workerCtx != nil {
		respondWithJSON(w, http.StatusOK, map[string]string{"status": "worker already running"})
		return
	}
	workerCtx, workerCancel = context.WithCancel(context.Background())
	go cfg.runFeedWorker(workerCtx, params)

	respondWithJSON(w, http.StatusOK, map[string]string{"status": "worker successfully started"})
}

func (cfg *apiConfig) stopFeedWorker(w http.ResponseWriter, r *http.Request) {
	if workerCancel == nil {
		respondWithJSON(w, http.StatusOK, map[string]string{"status": "Worker not running"})
		return
	}

	workerCancel()
	workerCtx = nil
	workerCancel = nil

	respondWithJSON(w, http.StatusOK, map[string]string{"status": "Feed worker stopped"})
}

func (cfg *apiConfig) runFeedWorker(ctx context.Context, params config) {
	ticker := time.NewTicker(time.Duration(params.Interval) * time.Second)
	defer ticker.Stop()

	log.Printf("Feed worker started.  Fetching every %d seconds with batch size %d\n", params.Interval, params.BatchSize)

	cfg.fetchAndProcessFeeds(ctx, params.BatchSize)

	for {
		select {
		case <-ctx.Done():
			log.Println("Feed Worker Stopped")
			return
		case <-ticker.C:
			cfg.fetchAndProcessFeeds(ctx, params.BatchSize)
		}
	}
}

func (cfg *apiConfig) fetchAndProcessFeeds(ctx context.Context, batchSize int) {
	log.Printf("Fetching next %d feeds\n", batchSize)

	feeds, err := cfg.DB.GetNextFeedToFetch(ctx, int32(batchSize))
	if err != nil {
		log.Printf("Failed to get feeds to fetch from database: %s\n", err)
		return
	}

	feedList := batchDatabaseFeedToFeeds(feeds)

	var wg sync.WaitGroup
	for _, feed := range feedList {
		wg.Add(1)
		go func(feed Feed) {
			defer wg.Done()
			processFeed(ctx, feed)
		}(feed)
	}
}

func processFeed(ctx context.Context, feed Feed) {
	client := &http.Client{
		Timeout: 10 * time.Second,
	}
	req, err := http.NewRequestWithContext(ctx, "GET", feed.Url, nil)
	if err != nil {
		log.Printf("Error creating request for %s: %v\n", feed.Url, err)
		return
	}

	resp, err := client.Do(req)
	if err != nil {
		log.Printf("Error fetching the RSS for %s: %v\n", feed.Url, err)
		return
	}
	defer resp.Body.Close()
	defer resp.Body.Close()

	rss := RSS{}
	xmlDecoder := xml.NewDecoder(resp.Body)
	err = xmlDecoder.Decode(&rss)
	if err != nil {
		log.Printf("Failed decode xml at: %s", feed.Url)
	}

	log.Printf("Successfully fetched and decoded RSS for %s\n", feed.Url)
	for _, item := range rss.Channel.Items {
		log.Printf("Source: %s, Title: %s", feed.FeedName, item.Title)
	}
}
