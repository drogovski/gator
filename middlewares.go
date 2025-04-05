package main

import (
	"context"
	"fmt"

	"github.com/drogovski/gator/internal/database"
)

func middlewareLoggedIn(handler func(s *state, cmd command, user database.User) error) func(*state, command) error {
	return func(s *state, cmd command) error {
		q := database.New(s.db)
		user, err := q.GetUser(context.Background(), s.cfg.CurrentUserName)
		if err != nil {
			return fmt.Errorf("cannot get logged in user: %w", err)
		}
		return handler(s, cmd, user)
	}
}
