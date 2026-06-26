package main

import (
	"context"
	"fmt"
	"strconv"

	"github.com/phpires/gator/internal/database"
)

func handlerBrowse(s *state, cmd command, user database.User) error {
	var limit int32
	limit = 2
	if len(cmd.Args) == 1 {
		userLimit, err := strconv.Atoi(cmd.Args[0])
		if err != nil {
			return fmt.Errorf("Error parsing argument %s\n", cmd.Args[0])
		}
		limit = int32(userLimit)
	}

	postForUser, err := s.dbState.GetPostsForUser(context.Background(), database.GetPostsForUserParams{
		UserID: user.ID,
		Limit:  limit,
	})
	if err != nil {
		return fmt.Errorf("Error finding posts for user %s\n", user.Name)
	}
	fmt.Printf("Found %d posts for user %s:\n", len(postForUser), user.Name)
	for _, post := range postForUser {
		fmt.Printf("%s from %s\n", post.PublishedAt.Time.Format("Mon Jan 2"), post.FeedName)
		fmt.Printf("--- %s ---\n", post.Title)
		fmt.Printf("    %v\n", post.Description.String)
		fmt.Printf("Link: %s\n", post.Url)
		fmt.Println("=====================================")
	}
	return nil
}
