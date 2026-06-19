package main

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/phpires/gator/internal/database"
)

func handlerFetch(s *state, cmd command) error {
	rssFeed, err := fetchFeed(context.Background(), "https://www.wagslane.dev/index.xml")
	if err != nil {
		return fmt.Errorf("couldn't fetch feed: %w", err)
	}
	fmt.Println(rssFeed)
	return nil
}

func handlerAddFeeds(s *state, cmd command) error {
	if len(cmd.Args) != 2 {
		return fmt.Errorf("usage: %v <name> <url>", cmd.Name)
	}
	feedName := cmd.Args[0]
	feedUrl := cmd.Args[1]

	currentUser, err := s.dbState.GetUserByName(context.Background(), s.configState.CurrentUserName)
	if err != nil {
		return fmt.Errorf("Error retrieving user: %w", err)
	}

	feed, err := s.dbState.CreateFeed(context.Background(), database.CreateFeedParams{
		ID:        uuid.New(),
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		Name:      feedName,
		Url:       feedUrl,
		UserID:    currentUser.ID,
	})
	if err != nil {
		return fmt.Errorf("Error creating feed: %w", err)
	}

	printFeed(feed)
	return nil
}

func printFeed(feed database.Feed) {
	fmt.Println("Feed created:")
	fmt.Printf(" * ID: %v\n", feed.ID)
	fmt.Printf(" * CreatedAt: %v\n", feed.CreatedAt)
	fmt.Printf(" * UpdatedAt: %v\n", feed.UpdatedAt)
	fmt.Printf(" * Name: %v\n", feed.Name)
	fmt.Printf(" * Url: %v\n", feed.Url)
	fmt.Printf(" * UserID: %v\n", feed.UserID)
	fmt.Println()
	fmt.Println("=====================================")
}
