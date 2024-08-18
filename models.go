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
	CreatedAt     time.Time
	UpdatedAt     time.Time
	LastFetchedAt time.Time
	FeedName      string
	Url           string
	ID            uuid.UUID
	UserID        uuid.UUID
}

func databaseFeedToFeed(feed database.Feed) Feed {
	return Feed{
		feed.CreatedAt,
		feed.UpdatedAt,
		feed.LastFetchedAt.Time,
		feed.FeedName,
		feed.Url,
		feed.ID,
		feed.UserID,
	}
}

func batchDatabaseFeedToFeeds(feed_lst []database.Feed) []Feed {
	output := make([]Feed, 0, len(feed_lst))
	for i := 0; i < len(feed_lst); i++ {
		output = append(output, databaseFeedToFeed(feed_lst[i]))
	}
	return output
}

type FeedFollow struct {
	CreatedAt time.Time
	UpdatedAt time.Time
	ID        uuid.UUID
	FeedID    uuid.UUID
	UserID    uuid.UUID
}

func databaseFeedFollowToFeedFollow(feedFollow database.FeedFollow) FeedFollow {
	return FeedFollow{
		CreatedAt: feedFollow.CreatedAt,
		UpdatedAt: feedFollow.UpdatedAt,
		ID:        feedFollow.ID,
		FeedID:    feedFollow.FeedID,
		UserID:    feedFollow.UserID,
	}
}

func selectAllFeedFollows(feedFollowLst []database.FeedFollow) []FeedFollow {
	output := make([]FeedFollow, 0, len(feedFollowLst))
	for i := 0; i < len(feedFollowLst); i++ {
		output = append(output, databaseFeedFollowToFeedFollow(feedFollowLst[i]))
	}
	return output
}
