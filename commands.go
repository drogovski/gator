package main

import (
	"errors"
)

type command struct {
	Name string
	Args []string
}

type commands struct {
	registeredCommands map[string]func(*state, command) error
}

func (c *commands) register(name string, f func(*state, command) error) {
	c.registeredCommands[name] = f
}

func (c *commands) run(s *state, cmd command) error {
	commandToRun, exists := c.registeredCommands[cmd.Name]
	if !exists {
		return errors.New("command with that name does not exist")
	}
	return commandToRun(s, cmd)
}
