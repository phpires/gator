package main

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"log"
	"time"

	"github.com/google/uuid"
	"github.com/phpires/gator/internal/database"
)

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
	username := cmd.args[0]
	_, err := s.dbState.GetUserByName(context.Background(), username)
	if err != nil {
		return err
	}

	err = s.configState.SetUser(username)
	if err != nil {
		return err
	}

	fmt.Println("User set.")
	return nil
}

func handlerRegister(s *state, cmd command) error {
	if len(cmd.args) == 0 {
		return errors.New("Name required to register.")
	}
	username := cmd.args[0]
	fmt.Printf("Registering username: %v\n", username)
	_, err := s.dbState.GetUserByName(context.Background(), username)
	if err != nil && err != sql.ErrNoRows {
		return err
	}

	if err == nil {
		log.Fatalf("Error: %v\n", err)
	}

	createUsr := database.CreateUserParams{
		ID:        uuid.New(),
		Name:      cmd.args[0],
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	user, err := s.dbState.CreateUser(context.Background(), createUsr)
	if err != nil {
		return err
	}

	log.Println("User created.")
	log.Printf("User(%v)\n", user)

	s.configState.SetUser(user.Name)
	return nil
}
