package app

import (
	"context"
	"fmt"
	"time"

	"github.com/Mikoofficial404/blog-go/internal/database"
	"github.com/google/uuid"
)

func handlerLogin(state *state, cmd command) error {
	if len(cmd.args) < 1 {
		return fmt.Errorf("username is required")
	}

	name := cmd.args[0]
	ctx := context.Background()
	user, err := state.db.GetUser(ctx, name)
	if err != nil {
		return fmt.Errorf("failed to get user %q: %w", name, err)
	}

	if err := state.cfg.SetUser(user.Name); err != nil {
		return fmt.Errorf("failed to set current user: %w", err)
	}

	fmt.Printf("Logged in as %s\n", name)
	return nil
}

func handlerUsers(state *state, _ command) error {
	ctx := context.Background()
	users, err := state.db.GetUsers(ctx)
	if err != nil {
		return fmt.Errorf("failed to list users: %w", err)
	}

	for _, user := range users {
		if user.Name == state.cfg.CurrentUserName {
			fmt.Printf("* %s (current)\n", user.Name)
		} else {
			fmt.Printf("* %s\n", user.Name)
		}
	}

	return nil
}

func handlerReset(state *state, _ command) error {
	ctx := context.Background()
	if err := state.db.DeleteUser(ctx); err != nil {
		return fmt.Errorf("failed to reset users: %w", err)
	}

	fmt.Println("Database reset")
	return nil
}

func handlerRegister(state *state, cmd command) error {
	if len(cmd.args) < 1 {
		return fmt.Errorf("username is required")
	}

	name := cmd.args[0]
	now := time.Now()
	params := database.CreateUserParams{
		ID:        uuid.New(),
		CreatedAt: now,
		UpdatedAt: now,
		Name:      name,
	}

	user, err := state.db.CreateUser(context.Background(), params)
	if err != nil {
		return fmt.Errorf("failed to create user %q: %w", name, err)
	}

	if err := state.cfg.SetUser(user.Name); err != nil {
		return fmt.Errorf("failed to set current user: %w", err)
	}

	fmt.Printf("Registered as %s\n", user.Name)
	return nil
}

func middlewareLoggedIn(handler func(*state, command, database.User) error) func(*state, command) error {
	return func(s *state, cmd command) error {
		if s.cfg.CurrentUserName == "" {
			return fmt.Errorf("no user is currently logged in")
		}

		user, err := s.db.GetUser(context.Background(), s.cfg.CurrentUserName)
		if err != nil {
			return fmt.Errorf("failed to get current user %q: %w", s.cfg.CurrentUserName, err)
		}

		return handler(s, cmd, user)
	}
}
