package main

import (
	"context"
	"fmt"
	"time"

	"github.com/drogovski/gator/internal/database"
)

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

	if len(feedFollows) == 0 {
		fmt.Println("No feed follows found for this user.")
		return nil
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

	fmt.Printf("%s unfollowed successfully!\n", feed.Name)
	return nil
}
