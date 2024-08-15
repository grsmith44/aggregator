package main

import (
	"time"

	"github.com/google/uuid"
	"github.com/grsmith44/aggregator.git/internal/database"
)

type User struct {
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	UserName  string    `json:"user_name"`
	APIKey    string    `json:"api_key"`
	ID        uuid.UUID `json:"id"`
}

func databaseUserToUser(user database.User) User {
	return User{
		user.CreatedAt,
		user.UpdatedAt,
		user.UserName,
		user.ApiKey,
		user.ID,
	}
}

type Feed struct {
	CreatedAt time.Time
	UpdatedAt time.Time
	FeedName  string
	Url       string
	ID        uuid.UUID
	UserID    uuid.UUID
}

func databaseFeedToFeed(feed database.Feed) Feed {
	return Feed{
		feed.CreatedAt,
		feed.UpdatedAt,
		feed.FeedName,
		feed.Url,
		feed.ID,
		feed.UserID,
	}
}

func selectAllfeeds(feed_lst []database.Feed) []Feed {
	output := make([]Feed, 0, len(feed_lst))
	for i := 0; i < len(feed_lst); i++ {
		output = append(output, databaseFeedToFeed(feed_lst[i]))
	}
	return output
}
