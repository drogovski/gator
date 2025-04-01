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
		log.Fatalf("User does not exist: %v\n", err)
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
	context := context.Background()

	user, err := s.db.CreateUser(context, database.CreateUserParams{
		ID:        uuid.New(),
		CreatedAt: time.Now().UTC(),
		UpdatedAt: time.Now().UTC(),
		Name:      username,
	})

	if err != nil {
		log.Fatalf("database error: %v", err)
	}

	err = s.cfg.SetUser(user.Name)
	if err != nil {
		log.Fatalf("cannot set user: %v", err)
	}
	fmt.Printf("The user with name %s was created.\n", user.Name)
	log.Default().Printf("User with this data was created: %v", user)
	return nil
}
