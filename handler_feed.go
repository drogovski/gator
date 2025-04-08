package main

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/drogovski/gator/internal/database"
	"github.com/drogovski/gator/internal/rss"
	"github.com/google/uuid"
)

func handlerAgg(s *state, cmd command) error {
	if len(cmd.Args) != 1 {
		return fmt.Errorf("usage: %s <time_between_reqs>", cmd.Name)
	}

	time_between_reqs, err := time.ParseDuration(cmd.Args[0])
	if err != nil {
		return fmt.Errorf("couldn't parse the duration from provided argument: %w",
			err)
	}

	fmt.Printf("Collecting feeds every %v\n", time_between_reqs)

	ticker := time.NewTicker(time_between_reqs)
	for ; ; <-ticker.C {
		err = scrapeFeeds(s)
		if err != nil {
			return err
		}
	}
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

func scrapeFeeds(s *state) error {
	q := database.New(s.db)
	feed, err := q.GetNextFeedToFetch(context.Background())
	if err != nil {
		return fmt.Errorf("couldn't get next feed to fetch: %w", err)
	}

	fetchedFeed, err := rss.FetchFeed(context.Background(), feed.Url)
	if err != nil {
		return fmt.Errorf("encountered error when trying to fetch feed: %w", err)
	}

	err = q.MarkFeedFetched(context.Background(), database.MarkFeedFetchedParams{
		ID: feed.ID,
		LastFetchedAt: sql.NullTime{
			Time:  time.Now().UTC(),
			Valid: true,
		},
	})

	if err != nil {
		return fmt.Errorf("couldn't update the last_fatched_at time: %w", err)
	}

	savePostsToDB(s, feed.ID, fetchedFeed)
	return nil
}

func savePostsToDB(s *state, feedId int32, fetchedFeed *rss.RSSFeed) {
	if len(fetchedFeed.Channel.Items) == 0 {
		fmt.Println("No new posts where fetched.")
		return
	}

	fmt.Printf("Saving %d posts to database...", len(fetchedFeed.Channel.Items))

	q := database.New(s.db)
	for _, item := range fetchedFeed.Channel.Items {
		pubDate, err := time.Parse(time.RFC1123Z, item.PubDate)
		if err != nil {
			fmt.Printf("couldn't parse the pubdate: %v\n", err)
		}

		_, err = q.CreatePost(context.Background(), database.CreatePostParams{
			ID:          uuid.New(),
			CreatedAt:   time.Now().UTC(),
			UpdatedAt:   time.Now().UTC(),
			Title:       item.Title,
			Url:         item.Link,
			Description: item.Description,
			PublishedAt: pubDate,
			FeedID:      feedId,
		})

		if err != nil {
			fmt.Printf("error when trying to save post: %v\n", err)
		}
	}
}
