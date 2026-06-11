package app

import (
	"context"
	"fmt"
	"time"

	"github.com/Mikoofficial404/blog-go/internal/database"
	"github.com/google/uuid"
)

func handlerAddfeed(state *state, cmd command, user database.User) error {
	if len(cmd.args) < 2 {
		return fmt.Errorf("feed name and URL are required")
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
		return fmt.Errorf("failed to create feed %q: %w", name, err)
	}

	fmt.Printf("Added feed: %s\n", feed.Name)

	_, err = state.db.CreateFeedFollow(context.Background(), database.CreateFeedFollowParams{
		UserID:    user.ID,
		FeedID:    feed.ID,
		ID:        uuid.New(),
		CreatedAt: now,
		UpdatedAt: now,
	})
	if err != nil {
		return fmt.Errorf("failed to follow feed %q: %w", feed.Name, err)
	}

	return nil
}

func handlerFeeds(state *state, _ command) error {
	feeds, err := state.db.JoinFeeds(context.Background())
	if err != nil {
		return fmt.Errorf("failed to list feeds: %w", err)
	}

	for _, item := range feeds {
		fmt.Printf("Feed name: %s\nURL: %s\nUser name: %s\n\n", item.Name, item.Url, item.UserName)
	}

	return nil
}

func handlerFollow(state *state, cmd command, user database.User) error {
	if len(cmd.args) < 1 {
		return fmt.Errorf("feed URL is required")
	}

	url := cmd.args[0]
	feed, err := state.db.GetFeedByUrl(context.Background(), url)
	if err != nil {
		return fmt.Errorf("failed to find feed by URL %q: %w", url, err)
	}

	now := time.Now()
	followed, err := state.db.CreateFeedFollow(context.Background(), database.CreateFeedFollowParams{
		ID:        uuid.New(),
		CreatedAt: now,
		UpdatedAt: now,
		UserID:    user.ID,
		FeedID:    feed.ID,
	})
	if err != nil {
		return fmt.Errorf("failed to follow feed %q: %w", feed.Name, err)
	}

	fmt.Printf("Followed feed: %s\n", followed.FeedName)
	fmt.Printf("Followed by: %s\n", followed.UserName)
	return nil
}

func handlerUnfollow(state *state, cmd command, user database.User) error {
	if len(cmd.args) < 1 {
		return fmt.Errorf("feed URL is required")
	}

	url := cmd.args[0]
	feed, err := state.db.GetFeedByUrl(context.Background(), url)
	if err != nil {
		return fmt.Errorf("failed to find feed by URL %q: %w", url, err)
	}

	if err := state.db.DeleteFeedFollow(context.Background(), database.DeleteFeedFollowParams{
		UserID: user.ID,
		FeedID: feed.ID,
	}); err != nil {
		return fmt.Errorf("failed to unfollow feed %q: %w", feed.Name, err)
	}

	fmt.Printf("Unfollowed feed: %s\n", feed.Name)
	return nil
}

func handlerFollowing(state *state, _ command, user database.User) error {
	feeds, err := state.db.GetFeedFollowsForUser(context.Background(), user.ID)
	if err != nil {
		return fmt.Errorf("failed to list followed feeds: %w", err)
	}

	for _, item := range feeds {
		fmt.Printf("Feed name: %s\n", item.FeedName)
	}

	return nil
}
