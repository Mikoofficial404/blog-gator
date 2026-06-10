package main

import (
	"context"
	"fmt"
	"time"

	"github.com/Mikoofficial404/blog-go/internal/database"
	"github.com/google/uuid"
)

func handlerAddfeed(state *state, cmd command, user database.User) error {
	if len(cmd.args) < 2 {
		return fmt.Errorf("Error ")
	}
	name := cmd.args[0]
	url := cmd.args[1]

	now := time.Now()
	feed, err := state.db.CreateFeeds(context.Background(), database.CreateFeedsParams{
		ID:        uuid.New(),
		CreatedAt: now,
		UpdatedAt: now,
		Name:      name,
		Url:       url,
		UserID:    user.ID,
	})
	if err != nil {
		return fmt.Errorf("Gagal")
	}
	fmt.Print(feed)

	_, err = state.db.CreateFeedFollow(context.Background(), database.CreateFeedFollowParams{
		UserID:    user.ID,
		FeedID:    feed.ID,
		ID:        uuid.New(),
		CreatedAt: now,
		UpdatedAt: now,
	})
	if err != nil {
		return fmt.Errorf("pesan: %w", err)
	}
	return nil
}
