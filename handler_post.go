package main

import (
	"context"
	"fmt"
	"strconv"

	"github.com/drogovski/gator/internal/database"
)

func handlerBrowse(s *state, cmd command, user database.User) error {
	if len(cmd.Args) != 1 {
		return fmt.Errorf("usage: %s <limit>", cmd.Name)
	}

	limit, err := strconv.ParseInt(cmd.Args[0], 10, 32)
	if err != nil {
		return fmt.Errorf("couldn't convert string to int: %w", err)
	}
	limit32 := int32(limit)

	q := database.New(s.db)
	posts, err := q.GetPostsForUser(context.Background(), database.GetPostsForUserParams{
		UserID: user.ID,
		Limit:  limit32,
	})
	if err != nil {
		return fmt.Errorf("couldn't get your posts from db: %w", err)
	}

	if len(posts) == 0 {
		fmt.Println("You have no posts to browse.")
		return nil
	}

	fmt.Println("Your posts to read:")
	for _, post := range posts {
		fmt.Println("**************************************")
		fmt.Printf("%s\n", post.Title)
		fmt.Println("**************************************")
		fmt.Printf("Pub Date: %s\n", post.PublishedAt)
		fmt.Printf("URL: %s\n", post.Url)
		fmt.Println()
		fmt.Printf("%s\n", post.Description)
		fmt.Println()
		fmt.Println()
	}
	return nil
}
