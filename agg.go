package main

import (
	"fmt"
	"time"
)

func handlerAgg(s *state, cmd command) error {
	if len(cmd.args) < 1 {
		return fmt.Errorf("not enough arguments were provided")
	}
	timeBetweenReqs := cmd.args[0]

	par, err := time.ParseDuration(timeBetweenReqs)
	if err != nil {
		return fmt.Errorf("error")
	}
	fmt.Printf("Collecting feeds every %s", par)
	ticker := time.NewTicker(par)
	for ; ; <-ticker.C {
		scrapeFeeds(s)
	}
}
