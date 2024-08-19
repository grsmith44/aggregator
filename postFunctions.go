package main

import (
	"context"
	"database/sql"
	"log"
	"time"

	"github.com/google/uuid"
	"github.com/grsmith44/aggregator.git/internal/database"
	"github.com/lib/pq"
)

func (cfg *apiConfig) addPostToDatabase(ctx context.Context, item Item, feed Feed) {
	publishedAt, err := time.Parse(time.RFC1123Z, item.PubDate)
	if err != nil {
		log.Printf("Failed to find Publish Date for: %s, link: %s", item.Title, item.Link)
		publishedAt = time.Time{}
	}
	_, err = cfg.DB.CreatePost(ctx, database.CreatePostParams{
		ID:          uuid.New(),
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
		Title:       item.Title,
		Url:         item.Link,
		Description: sql.NullString{String: item.Description, Valid: item.Description != ""},
		PublishedAt: sql.NullTime{Time: publishedAt, Valid: !publishedAt.IsZero()},
		FeedID:      feed.ID,
	})
	if err != nil {
		if isPgUniqueViolation(err) {
			return
		}
		log.Printf("error creating post: %s", err)
		return
	}
	log.Printf("Created new post entry: %s", item.Title)
}

func isPgUniqueViolation(err error) bool {
	if pqErr, ok := err.(*pq.Error); ok {
		return pqErr.Code == "23505"
	}
	return false
}
