package app

import (
	"context"
	"fmt"
	"strconv"

	"github.com/Mikoofficial404/blog-go/internal/database"
)

func handlerBrowse(s *state, cmd command, user database.User) error {
	limit := int32(2)

	if len(cmd.args) > 0 {
		parseLimit, err := strconv.Atoi(cmd.args[0])
		if err != nil {
			return fmt.Errorf("invalid limit: %w", err)
		}
		if parseLimit < 1 {
			return fmt.Errorf("limit must be greater than 0")
		}
		limit = int32(parseLimit)
	}

	posts, err := s.db.GetPostsForUser(context.Background(), database.GetPostsForUserParams{
		UserID: user.ID,
		Limit:  limit,
	})
	if err != nil {
		return fmt.Errorf("failed to load posts: %w", err)
	}

	for _, post := range posts {
		fmt.Printf("Title: %s\nURL: %s\n\n", post.Title, post.Url)
	}

	return nil
}
