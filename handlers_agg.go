package main

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"time"

	"github.com/google/uuid"
	"github.com/lib/pq"
	"github.com/lib/pq/pqerror"
	"github.com/phpires/gator/internal/database"
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
	if err != nil {
		return fmt.Errorf("Error marking feed %s as fetched: %w\n", nextFeed.Name, err)
	}

	feedName := nextFeed.Name
	url := nextFeed.Url

	rssFeed, err := fetchFeed(context.Background(), url)
	if err != nil {
		return fmt.Errorf("Error fetching feed %s %w\n", feedName, err)
	}

	for _, itemsFeed := range rssFeed.Channel.Item {
		fmt.Printf("Post: %s\n", itemsFeed.Title)
		parsedTime, err := parseTime(itemsFeed.PubDate)
		publishedAt := sql.NullTime{}
		if err == nil {
			publishedAt = sql.NullTime{
				Time:  parsedTime,
				Valid: true,
			}
		}

		_, err = s.dbState.CreatePost(context.Background(), database.CreatePostParams{
			ID:        uuid.New(),
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
			Title:     itemsFeed.Title,
			Url:       itemsFeed.Link,
			Description: sql.NullString{
				String: itemsFeed.Description,
				Valid:  true,
			},
			PublishedAt: publishedAt,
			FeedID:      nextFeed.ID,
		})

		if err != nil {
			if pqErr, ok := err.(*pq.Error); ok && pqErr.Code == pqerror.Code("23505") {
				continue
			}
			log.Printf("Couldn't create post: %v", err)
			continue
		}
	}
	log.Printf("Feed %s collected, %v posts found", feedName, len(rssFeed.Channel.Item))

	return nil
}
