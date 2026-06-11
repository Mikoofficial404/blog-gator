package main

import (
	"database/sql"
	"fmt"
	"os"

	"github.com/Mikoofficial404/blog-go/internal/app"
	"github.com/Mikoofficial404/blog-go/internal/config"
	"github.com/Mikoofficial404/blog-go/internal/database"
	_ "github.com/lib/pq"
)

func main() {
	cfg, err := config.Read()
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}

	args := os.Args[1:]
	if len(args) < 1 {
		fmt.Fprintln(os.Stderr, "at least one argument is required")
		os.Exit(1)
	}

	db, err := sql.Open("postgres", cfg.DbUrl)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	defer db.Close()

	queries := database.New(db)
	application := app.New(queries, &cfg)
	if err := application.Run(args); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
