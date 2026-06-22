package main

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/phpires/gator/internal/database"
)

func handlerFollowing(s *state, cmd command) error {
	currentUser, err := s.dbState.GetUserByName(context.Background(), s.configState.CurrentUserName)
	if err != nil {
		return fmt.Errorf("Error getting user by name: %w\n", err)
	}

	feedsFollowedByUser, err := s.dbState.GetFeedFollowsForUser(context.Background(), currentUser.ID)
	if err != nil {
		return fmt.Errorf("Error getting feed of user: %w\n", err)
	}

	fmt.Printf("Feed followed by user: %v", currentUser.ID)

	for _, feedsFollowed := range feedsFollowedByUser {
		printFeedFollowToUser(database.FeedFollow{
			ID:        feedsFollowed.ID,
			CreatedAt: feedsFollowed.CreatedAt,
			UpdatedAt: feedsFollowed.UpdatedAt,
			UserID:    feedsFollowed.UserID,
			FeedID:    feedsFollowed.FeedID,
		}, feedsFollowed.FeedName, feedsFollowed.UserName)
	}
	return nil
}

func handlerFeedFollow(s *state, cmd command) error {
	if len(cmd.Args) != 1 {
		return fmt.Errorf("usage: %v <url>", cmd.Name)
	}

	currentUser, err := s.dbState.GetUserByName(context.Background(), s.configState.CurrentUserName)
	if err != nil {
		return fmt.Errorf("Error getting user by name: %w\n", err)
	}

	feed, err := s.dbState.GetFeedByUrl(context.Background(), cmd.Args[0])

	createdFeedFollow, err := s.dbState.CreateFeedFollow(context.Background(), database.CreateFeedFollowParams{
		ID:        uuid.New(),
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		UserID:    currentUser.ID,
		FeedID:    feed.ID,
	})
	printFeedFollowToUser(database.FeedFollow{
		ID:        createdFeedFollow.ID,
		CreatedAt: createdFeedFollow.CreatedAt,
		UpdatedAt: createdFeedFollow.UpdatedAt,
		UserID:    createdFeedFollow.UserID,
		FeedID:    createdFeedFollow.FeedID,
	}, createdFeedFollow.FeedName, createdFeedFollow.UserName)
	return nil
}

func printFeedFollowToUser(feedFollow database.FeedFollow, feedName string, userName string) {
	fmt.Println("Feed followed:")
	fmt.Printf(" * ID: %v\n", feedFollow.ID)
	fmt.Printf(" * CreatedAt: %v\n", feedFollow.CreatedAt)
	fmt.Printf(" * UpdatedAt: %v\n", feedFollow.UpdatedAt)
	fmt.Printf(" * UserID: %v\n", feedFollow.UserID)
	fmt.Printf(" * FeedID: %v\n", feedFollow.FeedID)
	fmt.Printf(" * FeedName: %v\n", feedName)
	fmt.Printf(" * UserName: %v\n", userName)
	fmt.Println()
	fmt.Println("=====================================")
}
