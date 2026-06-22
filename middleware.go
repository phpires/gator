package main

import (
	"context"
	"fmt"

	"github.com/phpires/gator/internal/database"
)

func middlewareLoggedIn(handler func(s *state, cmd command, user database.User) error) func(*state, command) error {
	return func(s *state, cmd command) error {
		currentUser, err := s.dbState.GetUserByName(context.Background(), s.configState.CurrentUserName)
		if err != nil {
			return fmt.Errorf("Error getting logged in user: %w\n", err)
		}

		err = handler(s, cmd, currentUser)
		if err != nil {
			return fmt.Errorf("Error calling cmd '%s' handler: %w\n", cmd.Name, err)
		}
		return nil
	}
}
