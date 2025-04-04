package main

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/drogovski/gator/internal/database"
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

func handlerAddFeed(s *state, cmd command) error {
	if len(cmd.Args) != 2 {
		log.Fatal("You have to provide feed name and URL as arguments")
	}
	currentUserName := s.cfg.CurrentUserName
	currentUser, err := s.db.GetUser(context.Background(), currentUserName)
	if err != nil {
		log.Fatalf("Error when getting current user: %v", err)
	}

	feedName := cmd.Args[0]
	url := cmd.Args[1]

	feed, err := s.db.CreateFeed(context.Background(), database.CreateFeedParams{
		Name:      feedName,
		Url:       url,
		UserID:    currentUser.ID,
		CreatedAt: time.Now().UTC(),
		UpdatedAt: time.Now().UTC(),
	})

	if err != nil {
		log.Fatalf("Error when trying to create new feed: %v", feed)
	}

	fmt.Printf("Created new feed:\n %v\n", feed)
	return nil
}

func handlerFeeds(s *state, cmd command) error {
	feeds, err := s.db.GetFeeds(context.Background())
	if err != nil {
		log.Fatalf("Couldn't get feeds from db: %v", err)
	}
	printFeeds(feeds)
	return nil
}

func printFeeds(feeds []database.GetFeedsRow) {
	fmt.Println("Your Feeds:")
	for _, feed := range feeds {
		fmt.Printf(" * %s | url: %s | author: %s\n", feed.Name, feed.Url, feed.Name_2)
	}
}
