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

func handlerAddFeed(s *state, cmd command, user database.User) error {
	if len(cmd.Args) != 2 {
		return fmt.Errorf("usage: %s <name> <url>", cmd.Name)
	}

	tx, err := s.db.Begin()
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer tx.Rollback()

	qtx := database.New(tx)

	feedName := cmd.Args[0]
	url := cmd.Args[1]

	feed, err := qtx.CreateFeed(context.Background(), database.CreateFeedParams{
		Name:      feedName,
		Url:       url,
		UserID:    user.ID,
		CreatedAt: time.Now().UTC(),
		UpdatedAt: time.Now().UTC(),
	})
	if err != nil {
		return fmt.Errorf("error when trying to create new feed: %v", err)
	}

	_, err = qtx.CreateFeedFollow(context.Background(), database.CreateFeedFollowParams{
		CreatedAt: time.Now().UTC(),
		UpdatedAt: time.Now().UTC(),
		UserID:    user.ID,
		FeedID:    feed.ID,
	})
	if err != nil {
		return fmt.Errorf("couldn't create feed follow: %w", err)
	}

	if err := tx.Commit(); err != nil {
		return fmt.Errorf("failed to commit transaction: %w", err)
	}

	fmt.Printf("Created new feed:\n %v\n", feed)
	return nil
}

func handlerFeeds(s *state, cmd command) error {
	q := database.New(s.db)
	feeds, err := q.GetFeeds(context.Background())
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

func handlerFollow(s *state, cmd command, user database.User) error {
	if len(cmd.Args) != 1 {
		return fmt.Errorf("usage: %s <feed_url>", cmd.Name)
	}

	feedURL := cmd.Args[0]
	q := database.New(s.db)
	feed, err := q.GetFeedByUrl(context.Background(), feedURL)

	if err != nil {
		return fmt.Errorf("couldn't get feed with given url: %w", err)
	}

	feedFollow, err := q.CreateFeedFollow(context.Background(), database.CreateFeedFollowParams{
		CreatedAt: time.Now().UTC(),
		UpdatedAt: time.Now().UTC(),
		UserID:    user.ID,
		FeedID:    feed.ID,
	})

	if err != nil {
		return fmt.Errorf("couldn't create feed follow: %w", err)
	}

	fmt.Printf("Created %s feed follow by %s.", feedFollow.FeedName, feedFollow.UserName)
	return nil
}

func handlerFollowing(s *state, cmd command, user database.User) error {
	q := database.New(s.db)
	feedFollows, err := q.GetFeedFollowsForUser(context.Background(), user.Name)
	if err != nil {
		return fmt.Errorf("couldn't get following information for current user: %w", err)
	}

	fmt.Printf("Feeds followed by: %s\n", s.cfg.CurrentUserName)
	fmt.Println("==========================================")
	for _, feed := range feedFollows {
		fmt.Printf(" * %s\n", feed.FeedName)
	}
	return nil
}

func handlerUnfollow(s *state, cmd command, user database.User) error {
	if len(cmd.Args) != 1 {
		return fmt.Errorf("usage: %s <feed_url>", cmd.Name)
	}
	feedURL := cmd.Args[0]

	q := database.New(s.db)
	feed, err := q.GetFeedByUrl(context.Background(), feedURL)
	if err != nil {
		return fmt.Errorf("couldn't get feed for given url: %w", err)
	}

	err = q.DeleteFeedFollow(context.Background(), database.DeleteFeedFollowParams{
		FeedID: feed.ID,
		UserID: user.ID,
	})

	if err != nil {
		return fmt.Errorf("couldn't unfollow: %w", err)
	}

	return nil
}
