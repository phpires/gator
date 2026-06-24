package main

import (
	"context"
	"fmt"
	"log"
	"time"
)

func handlerAgg(s *state, cmd command) error {
	if len(cmd.Args) != 1 {
		return fmt.Errorf("usage: %v <time_between_reqs>", cmd.Name)
	}

	timeBetweenReqs, err := time.ParseDuration(cmd.Args[0])
	if err != nil {
		return fmt.Errorf("Error parsing duration %w", err)
	}

	log.Printf("Collecting feeds every %s...", timeBetweenReqs)
	log.Printf("Kill with ctrl + c...")
	ticker := time.NewTicker(timeBetweenReqs)
	for ; ; <-ticker.C {
		scrapeFeeds(s)
	}

}

func scrapeFeeds(s *state) error {
	nextFeed, err := s.dbState.GetNextFeedToFetch(context.Background())
	if err != nil {
		return fmt.Errorf("Error getting next feed to fetch %w \n", err)
	}

	nextFeed, err = s.dbState.MarkFeedFetched(context.Background(), nextFeed.ID)
	feedName := nextFeed.Name
	if err != nil {
		return fmt.Errorf("Error marking feed %s as fetched: %w\n", feedName, err)
	}

	rssFeed, err := fetchFeed(context.Background(), nextFeed.Url)
	if err != nil {
		return fmt.Errorf("Error fetching feed %s %w\n", feedName, err)
	}

	for _, itemsFeed := range rssFeed.Channel.Item {
		fmt.Printf("Post: %s\n", itemsFeed.Title)
	}
	log.Printf("Feed %s collected, %v posts found", feedName, len(rssFeed.Channel.Item))

	return nil
}
