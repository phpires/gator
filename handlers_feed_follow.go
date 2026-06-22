package main

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/phpires/gator/internal/database"
)

func handlerFollowing(s *state, cmd command, user database.User) error {

	feedsFollowedByUser, err := s.dbState.GetFeedFollowsForUser(context.Background(), user.ID)
	if err != nil {
		return fmt.Errorf("Error getting feed of user: %w\n", err)
	}

	fmt.Printf("Feed followed by user: %v", user.ID)

	for _, feedsFollowed := range feedsFollowedByUser {
		printFeedFollowToUser(feedsFollowed.FeedName, feedsFollowed.UserName)
	}
	return nil
}

func handlerFeedFollow(s *state, cmd command, user database.User) error {
	if len(cmd.Args) != 1 {
		return fmt.Errorf("usage: %v <feed_url>", cmd.Name)
	}

	feed, err := s.dbState.GetFeedByUrl(context.Background(), cmd.Args[0])
	if err != nil {
		return fmt.Errorf("Error getting feed by url: %w", err)
	}

	createdFeedFollow, err := s.dbState.CreateFeedFollow(context.Background(), database.CreateFeedFollowParams{
		ID:        uuid.New(),
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		UserID:    user.ID,
		FeedID:    feed.ID,
	})
	fmt.Println("Feed followed created:")
	printFeedFollowToUser(createdFeedFollow.FeedName, createdFeedFollow.UserName)
	return nil
}

func printFeedFollowToUser(feedName string, userName string) {
	fmt.Printf(" * User: 	%s\n", userName)
	fmt.Printf(" * Feed: 	%s\n", feedName)
	fmt.Println()
	fmt.Println("=====================================")
}
