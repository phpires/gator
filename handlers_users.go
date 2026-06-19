package main

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/google/uuid"
	"github.com/phpires/gator/internal/database"
)

func handlerLogin(s *state, cmd command) error {
	if len(cmd.Args) == 0 {
		return fmt.Errorf("usage: %v <name>", cmd.Name)
	}
	username := cmd.Args[0]
	_, err := s.dbState.GetUserByName(context.Background(), username)
	if err != nil {
		return fmt.Errorf("Couldn't find user: %w", err)
	}

	err = s.configState.SetUser(username)
	if err != nil {
		return fmt.Errorf("Couldn't set current user: %w", err)
	}

	fmt.Println("User set.")
	return nil
}

func handlerRegister(s *state, cmd command) error {
	if len(cmd.Args) == 0 {
		return fmt.Errorf("usage: %v <name>", cmd.Name)
	}

	createUsr := database.CreateUserParams{
		ID:        uuid.New(),
		Name:      cmd.Args[0],
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	user, err := s.dbState.CreateUser(context.Background(), createUsr)
	if err != nil {
		return fmt.Errorf("Error creating user: %w", err)
	}

	log.Println("User created.")
	log.Printf("User(%v)\n", user)

	err = s.configState.SetUser(user.Name)
	if err != nil {
		return fmt.Errorf("Error setting user: %w", err)
	}
	return nil
}

func handlerReset(s *state, cmd command) error {
	err := s.dbState.DeleteAllUsers(context.Background())
	if err != nil {
		return err
	}
	log.Println("Users removed.")
	s.configState.SetUser("")
	return nil
}
