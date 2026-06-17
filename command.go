package main

import (
	"errors"
	"fmt"

	"github.com/phpires/gator/internal/config"
)

type state struct {
	configState *config.Config
}

type command struct {
	name string
	args []string
}

type commands struct {
	handlers map[string]func(*state, command) error
}

func (c *commands) run(s *state, cmd command) error {
	commandFunction, ok := c.handlers[cmd.name]
	if !ok {
		return errors.New("No function registered.")
	}
	return commandFunction(s, cmd)
}

func (c *commands) register(name string, f func(*state, command) error) {
	c.handlers[name] = f
}

func handlerLogin(s *state, cmd command) error {
	if len(cmd.args) == 0 {
		return errors.New("Username is required.")
	}

	err := s.configState.SetUser(cmd.args[0])
	if err != nil {
		return err
	}

	fmt.Println("User set.")
	return nil
}
