package main

import (
	"errors"

	"github.com/drogovski/gator/internal/config"
)

type state struct {
	cfg *config.Config
}

type command struct {
	name string
	args []string
}

type commands struct {
	cmds map[string]func(*state, command) error
}

func (c *commands) register(name string, f func(*state, command) error) {
	c.cmds[name] = f
}

func (c *commands) run(s *state, cmd command) error {
	commandToRun, exists := c.cmds[cmd.name]
	if !exists {
		return errors.New("command with that name does not exist")
	}
	return commandToRun(s, cmd)
}
