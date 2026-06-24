package main

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/phpires/gator/internal/database"
)

func handlerListFeeds(s *state, cmd command) error {
	feeds, err := s.dbState.ListFeeds(context.Background())
	if err != nil {
		return fmt.Errorf("Couldn't list feed: %w", err)
	}

	if len(feeds) == 0 {
		fmt.Println("No feeds found.")
		return nil
	}

	var userCreator database.User
	for _, feed := range feeds {
		userCreator, err = s.dbState.GetUserById(context.Background(), feed.UserID)
		if err != nil {
			return fmt.Errorf("Error getting user that created feed: %w", err)
		}
		printFeedToUser(feed, userCreator)
	}
	return nil
}

func handlerAddFeeds(s *state, cmd command, user database.User) error {
	if len(cmd.Args) != 2 {
		return fmt.Errorf("usage: %v <name> <url>", cmd.Name)
	}
	feedName := cmd.Args[0]
	feedUrl := cmd.Args[1]

	feed, err := s.dbState.CreateFeed(context.Background(), database.CreateFeedParams{
		ID:        uuid.New(),
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		Name:      feedName,
		Url:       feedUrl,
		UserID:    user.ID,
	})
	if err != nil {
		return fmt.Errorf("Error creating feed: %w", err)
	}

	fmt.Println("Linking feed to user.")
	_, err = s.dbState.CreateFeedFollow(context.Background(), database.CreateFeedFollowParams{
		ID:        uuid.New(),
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		UserID:    user.ID,
		FeedID:    feed.ID,
	})
	if err != nil {
		return fmt.Errorf("Error linking feed to user: %w", err)
	}

	printFeedToUser(feed, user)
	return nil
}

func printFeedToUser(feed database.Feed, user database.User) {
	fmt.Println("Feed created:")
	fmt.Printf(" * ID: %v\n", feed.ID)
	fmt.Printf(" * CreatedAt: %v\n", feed.CreatedAt)
	fmt.Printf(" * UpdatedAt: %v\n", feed.UpdatedAt)
	fmt.Printf(" * Name: %v\n", feed.Name)
	fmt.Printf(" * Url: %v\n", feed.Url)
	fmt.Printf(" * User: %v\n", user.Name)
	fmt.Println()
	fmt.Println("=====================================")
}
