package main

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"os"
	"strconv"
	"time"

	"github.com/Mikoofficial404/blog-go/internal/config"
	"github.com/Mikoofficial404/blog-go/internal/database"
	"github.com/google/uuid"
	_ "github.com/lib/pq"
)

type state struct {
	db  *database.Queries
	cfg *config.Config
}
type command struct {
	name string
	args []string
}
type commands struct {
	handlers map[string]func(*state, command) error
}

func handlerLogin(state *state, cmd command) error {
	if len(cmd.args) < 1 {
		return fmt.Errorf("username is required")
	}
	name := cmd.args[0]
	user, err := state.db.GetUser(context.Background(), name)
	if err != nil {
		return err
	}
	err = state.cfg.SetUser(user.Name)
	if err != nil {
		return err
	}

	fmt.Println("logged in as", name)
	return nil
}

func handleUser(state *state, cmd command) error {
	value, err := state.db.GetUsers(context.Background())
	if err != nil {
		return err
	}
	for _, user := range value {
		if user.Name == state.cfg.CurrentUserName {
			fmt.Printf("* %s (current)\n", user.Name)
		} else {
			fmt.Printf("* %s\n", user.Name)
		}

	}
	return nil
}

func handlerReset(state *state, cmd command) error {
	err := state.db.DeleteUser(context.Background())
	if err != nil {
		return err
	}
	fmt.Println("reset")
	return nil
}

func handlerFollow(state *state, cmd command, user database.User) error {
	if len(cmd.args) < 1 {
		return fmt.Errorf("username is required")
	}
	url := cmd.args[0]

	feed, err := state.db.GetFeedByUrl(context.Background(), url)
	if err != nil {
		return fmt.Errorf("pesan: %w", err)
	}
	now := time.Now()
	createFeed, err := state.db.CreateFeedFollow(context.Background(), database.CreateFeedFollowParams{
		ID:        uuid.New(),
		CreatedAt: now,
		UpdatedAt: now,
		UserID:    user.ID,
		FeedID:    feed.ID,
	})
	if err != nil {
		return fmt.Errorf("pesan: %w", err)
	}
	fmt.Println("follow", createFeed.FeedName)
	fmt.Println("followbyName", createFeed.UserName)
	return nil
}

func handlerUnfollow(state *state, cmd command, user database.User) error {
	if len(cmd.args) < 1 {
		return fmt.Errorf("not enough arguments were provided")
	}
	url := cmd.args[0]
	feed, err := state.db.GetFeedByUrl(context.Background(), url)
	if err != nil {
		return fmt.Errorf("pesan: %w", err)
	}
	err = state.db.DeleteFeedFollow(context.Background(), database.DeleteFeedFollowParams{
		UserID: user.ID,
		FeedID: feed.ID,
	})
	if err != nil {
		return fmt.Errorf("pesan: %w", err)
	}
	return nil
}

func middlewareLoggedIn(handler func(s *state, cmd command, user database.User) error) func(*state, command) error {
	return func(s *state, cmd command) error {
		usr := s.cfg.CurrentUserName
		getUsr, err := s.db.GetUser(context.Background(), usr)
		if err != nil {
			return err
		}
		return handler(s, cmd, getUsr)
	}
}

func handlerFeeds(state *state, cmd command) error {
	value, err := state.db.JoinFeeds(context.Background())
	if err != nil {
		return err
	}
	for _, item := range value {
		fmt.Println("Nama feed:", item.Name)
		fmt.Println("URL:", item.Url)
		fmt.Println("UserName:", item.UserName)
	}
	return nil
}

func handlerFollowing(state *state, cmd command, user database.User) error {

	feed, err := state.db.GetFeedFollowsForUser(context.Background(), user.ID)
	if err != nil {
		return fmt.Errorf("pesan: %w", err)
	}
	for _, item := range feed {
		fmt.Println("Name", item.FeedName)
	}
	return nil
}

func handlerRegister(state *state, cmd command) error {
	if len(cmd.args) < 1 {
		return fmt.Errorf("not enough arguments were provided")
	}

	name := cmd.args[0]
	now := time.Now()
	datable := database.CreateUserParams{
		ID:        uuid.New(),
		CreatedAt: now,
		UpdatedAt: now,
		Name:      name,
	}
	db, err := state.db.CreateUser(context.Background(), datable)
	if err != nil {
		return err
	}
	err = state.cfg.SetUser(datable.Name)
	if err != nil {
		return err
	}
	fmt.Println("registered as", db.Name)
	return nil
}

func scrapeFeeds(s *state) error {
	nextFeed, err := s.db.GetNextFeedToFetch(context.Background())
	if err != nil {
		return fmt.Errorf("Error")
	}
	err = s.db.MarkFeedFetched(context.Background(), nextFeed.ID)
	if err != nil {
		return fmt.Errorf("Error")
	}
	value, err := fetchFeed(context.Background(), nextFeed.Url)
	if err != nil {
		return fmt.Errorf("Error")
	}
	now := time.Now()

	for _, item := range value.Channel.Item {
		pubDate, err := time.Parse(time.RFC1123Z, item.PubDate)
		if err != nil {
			log.Printf("error parsing date: %v", err)
			continue // skip item ini, lanjut ke item berikutnya
		}
		_, err = s.db.CreatePost(context.Background(), database.CreatePostParams{
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
			log.Printf("error saving post: %v", err)
		}
	}
	return nil
}

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
	postUsr, err := s.db.GetPostsForUser(context.Background(), database.GetPostsForUserParams{
		UserID: user.ID,
		Limit:  limit,
	})
	if err != nil {
		return fmt.Errorf("invalid")
	}
	for _, value := range postUsr {
		fmt.Printf("Title: %s\nURL: %s\n\n", value.Title, value.Url)
	}
	return nil
}
func newCommands() *commands {
	return &commands{
		handlers: make(map[string]func(*state, command) error),
	}
}

func (registry *commands) register(name string, handler func(*state, command) error) {
	if registry.handlers == nil {
		return
	}
	registry.handlers[name] = handler
}
func (registry *commands) run(st *state, cmd command) error {
	handler, exists := registry.handlers[cmd.name]
	if !exists {
		return fmt.Errorf("command %q tidak ditemukan", cmd.name)
	}

	return handler(st, cmd)
}

func main() {
	cr, err := config.Read()
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	args := os.Args[1:]
	if len(args) < 1 {
		fmt.Fprintln(os.Stderr, "not enough arguments were provided")
		os.Exit(1)
	}

	db, err := sql.Open("postgres", cr.DbUrl)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}

	dbQueries := database.New(db)

	defer db.Close()

	appState := state{
		db:  dbQueries,
		cfg: &cr,
	}
	cmd := command{
		name: args[0],
		args: args[1:],
	}

	registry := newCommands()
	registry.register("login", handlerLogin)
	registry.register("register", handlerRegister)
	registry.register("reset", handlerReset)
	registry.register("users", handleUser)
	registry.register("agg", handlerAgg)
	registry.register("addfeed", middlewareLoggedIn(handlerAddfeed))
	registry.register("feeds", handlerFeeds)
	registry.register("follow", middlewareLoggedIn(handlerFollow))
	registry.register("following", middlewareLoggedIn(handlerFollowing))
	registry.register("unfollow", middlewareLoggedIn(handlerUnfollow))
	registry.register("browse", middlewareLoggedIn(handlerBrowse))
	err = registry.run(&appState, cmd)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
