package app

import (
	"fmt"

	"github.com/Mikoofficial404/blog-go/internal/config"
	"github.com/Mikoofficial404/blog-go/internal/database"
)

type state struct {
	db  *database.Queries
	cfg *config.Config
}

type command struct {
	name string
	args []string
}

type App struct {
	state    *state
	commands *commands
}

func New(db *database.Queries, cfg *config.Config) *App {
	appState := &state{
		db:  db,
		cfg: cfg,
	}

	registry := newCommands()
	registry.register("login", handlerLogin)
	registry.register("register", handlerRegister)
	registry.register("reset", handlerReset)
	registry.register("users", handlerUsers)
	registry.register("agg", handlerAgg)
	registry.register("addfeed", middlewareLoggedIn(handlerAddfeed))
	registry.register("feeds", handlerFeeds)
	registry.register("follow", middlewareLoggedIn(handlerFollow))
	registry.register("following", middlewareLoggedIn(handlerFollowing))
	registry.register("unfollow", middlewareLoggedIn(handlerUnfollow))
	registry.register("browse", middlewareLoggedIn(handlerBrowse))

	return &App{
		state:    appState,
		commands: registry,
	}
}

func (a *App) Run(args []string) error {
	if len(args) < 1 {
		return fmt.Errorf("at least one argument is required")
	}

	cmd := command{
		name: args[0],
		args: args[1:],
	}

	return a.commands.run(a.state, cmd)
}
