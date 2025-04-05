package main

import (
	"context"
	"errors"
	"fmt"
	"log"
	"time"

	"github.com/drogovski/gator/internal/database"
	"github.com/google/uuid"
)

func handlerLogin(s *state, cmd command) error {
	if len(cmd.Args) != 1 {
		return errors.New("you have to provide username as argument")
	}
	username := cmd.Args[0]
	user, err := s.db.GetUser(context.Background(), username)
	if err != nil {
		return fmt.Errorf("user does not exist: %v", err)
	}

	err = s.cfg.SetUser(user.Name)
	if err != nil {
		return err
	}
	fmt.Printf("User <%s> has been set.\n", username)
	return nil
}

func handlerRegister(s *state, cmd command) error {
	if len(cmd.Args) != 1 {
		return errors.New("you have to provide username as argument")
	}
	username := cmd.Args[0]
	user, err := s.db.CreateUser(context.Background(), database.CreateUserParams{
		ID:        uuid.New(),
		CreatedAt: time.Now().UTC(),
		UpdatedAt: time.Now().UTC(),
		Name:      username,
	})

	if err != nil {
		return fmt.Errorf("database error: %v", err)
	}

	err = s.cfg.SetUser(user.Name)
	if err != nil {
		return fmt.Errorf("cannot set user: %v", err)
	}
	fmt.Printf("The user with name %s was created.\n", user.Name)
	log.Default().Printf("User with this data was created: %v", user)
	return nil
}

func handlerReset(s *state, cmd command) error {
	err := s.db.DeleteAllUsers(context.Background())

	if err != nil {
		return fmt.Errorf("cannot reset db: %v", err)
	}
	fmt.Println("Database was successfully reset.")
	return nil
}

func handlerUsers(s *state, cmd command) error {
	users, err := s.db.GetUsers(context.Background())
	if err != nil {
		return fmt.Errorf("couldn't get users from db: %v", users)
	}
	printUsers(s, users)
	return nil
}

func printUsers(s *state, users []database.User) {
	for _, user := range users {
		if user.Name == s.cfg.CurrentUserName {
			fmt.Printf(" * %s (current)\n", user.Name)
			continue
		}
		fmt.Printf(" * %s\n", user.Name)
	}
}
