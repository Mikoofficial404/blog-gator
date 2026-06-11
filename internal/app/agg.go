package app

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"time"

	"github.com/Mikoofficial404/blog-go/internal/database"
	"github.com/google/uuid"
)

func handlerAgg(s *state, cmd command) error {
	if len(cmd.args) < 1 {
		return fmt.Errorf("duration is required")
	}

	timeBetweenReqs := cmd.args[0]
	duration, err := time.ParseDuration(timeBetweenReqs)
	if err != nil {
		return fmt.Errorf("invalid duration: %w", err)
	}

	fmt.Printf("Collecting feeds every %s\n", duration)
	ticker := time.NewTicker(duration)
	defer ticker.Stop()

	for range ticker.C {
		_ = scrapeFeeds(s)
	}

	return nil
}

func scrapeFeeds(s *state) error {
	ctx := context.Background()
	nextFeed, err := s.db.GetNextFeedToFetch(ctx)
	if err != nil {
		return fmt.Errorf("failed to get next feed to fetch: %w", err)
	}

	if err := s.db.MarkFeedFetched(ctx, nextFeed.ID); err != nil {
		return fmt.Errorf("failed to mark feed %q as fetched: %w", nextFeed.Name, err)
	}

	value, err := fetchFeed(ctx, nextFeed.Url)
	if err != nil {
		return fmt.Errorf("failed to fetch feed %q: %w", nextFeed.Url, err)
	}

	now := time.Now()
	for _, item := range value.Channel.Item {
		pubDate, err := time.Parse(time.RFC1123Z, item.PubDate)
		if err != nil {
			log.Printf("failed to parse publication date for %q: %v", item.Title, err)
			continue
		}

		_, err = s.db.CreatePost(ctx, database.CreatePostParams{
			ID:          uuid.New(),
			CreatedAt:   now,
			UpdatedAt:   now,
			Title:       item.Title,
			Url:         item.Link,
			Description: sql.NullString{String: item.Description, Valid: item.Description != ""},
			PublishedAt: sql.NullString{String: pubDate.Format(time.RFC1123Z), Valid: true},
			FeedID:      nextFeed.ID,
		})
		if err != nil {
			log.Printf("failed to save post %q: %v", item.Title, err)
		}
	}

	return nil
}
