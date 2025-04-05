package main

import (
	"context"
	"fmt"
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
		return fmt.Errorf("encountered error when trying to fetch feed: %v", err)
	}
	fmt.Printf("%v\n", feed)
	return nil
}

func handlerAddFeed(s *state, cmd command) error {
	if len(cmd.Args) != 2 {
		return fmt.Errorf("usage: %s <name> <url>", cmd.Name)
	}
	currentUser, err := s.db.GetUser(context.Background(), s.cfg.CurrentUserName)
	if err != nil {
		return fmt.Errorf("error when getting current user: %v", err)
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
		return fmt.Errorf("error when trying to create new feed: %v", err)
	}

	fmt.Printf("Created new feed:\n %v\n", feed)
	return nil
}

func handlerFeeds(s *state, cmd command) error {
	feeds, err := s.db.GetFeeds(context.Background())
	if err != nil {
		return fmt.Errorf("couldn't get feeds from db: %v", err)
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

func handlerFollow(s *state, cmd command) error {
	if len(cmd.Args) != 1 {
		return fmt.Errorf("usage: %s <feed_url>", cmd.Name)
	}

	feedURL := cmd.Args[0]
	feed, err := s.db.GetFeedByUrl(context.Background(), feedURL)

	if err != nil {
		return fmt.Errorf("couldn't get feed with given url: %w", err)
	}

	currentUser, err := s.db.GetUser(context.Background(), s.cfg.CurrentUserName)
	if err != nil {
		return fmt.Errorf("couldn't get user with given url: %w", err)
	}

	feedFollow, err := s.db.CreateFeedFollow(context.Background(), database.CreateFeedFollowParams{
		CreatedAt: time.Now().UTC(),
		UpdatedAt: time.Now().UTC(),
		UserID:    currentUser.ID,
		FeedID:    feed.ID,
	})

	if err != nil {
		return fmt.Errorf("couldn't create feed follow: %w", err)
	}

	fmt.Printf("Created %s feed follow by %s.", feedFollow.FeedName, feedFollow.UserName)
	return nil
}
