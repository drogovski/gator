package main

import (
	"database/sql"
	"log"
	"os"

	"github.com/drogovski/gator/internal/config"
	_ "github.com/lib/pq"
)

type state struct {
	db  *sql.DB
	cfg *config.Config
}

func main() {
	cfg, err := config.Read()
	if err != nil {
		log.Fatalf("error reading config: %v", err)
	}

	db, err := sql.Open("postgres", cfg.DBURL)
	if err != nil {
		log.Fatalf("error creating connection to db: %v", err)
	}

	programState := &state{
		db:  db,
		cfg: &cfg,
	}

	cmds := commands{
		registeredCommands: map[string]func(*state, command) error{},
	}

	cmds.register("addfeed", middlewareLoggedIn(handlerAddFeed))
	cmds.register("agg", handlerAgg)
	cmds.register("feeds", handlerFeeds)
	cmds.register("follow", middlewareLoggedIn(handlerFollow))
	cmds.register("following", middlewareLoggedIn(handlerFollowing))
	cmds.register("login", handlerLogin)
	cmds.register("register", handlerRegister)
	cmds.register("reset", handlerReset)
	cmds.register("unfollow", middlewareLoggedIn(handlerUnfollow))
	cmds.register("users", handlerUsers)

	args := os.Args
	if len(args) < 2 {
		log.Fatal("Usage: cli <command> [args...]")
		return
	}

	cmd := command{
		Name: args[1],
		Args: args[2:],
	}

	err = cmds.run(programState, cmd)
	if err != nil {
		log.Fatal(err)
	}
}
