package main

import (
	"context"
	"fmt"
	"log"

	"github.com/drogovski/gator/internal/rss"
)

const (
	feedURL = "https://www.wagslane.dev/index.xml"
)

func handlerAgg(s *state, cmd command) error {
	feed, err := rss.FetchFeed(context.Background(), feedURL)
	if err != nil {
		log.Fatalf("Encountered error when trying to fetch feed: %v", err)
	}
	fmt.Printf("%v", feed)
	return nil
}
