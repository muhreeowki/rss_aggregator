package main

import (
	"context"
	"database/sql"
	"log"
	"strings"
	"sync"
	"time"

	"github.com/google/uuid"
	"github.com/muhreeowki/rss_aggregator/internal/database"
)

// startScraping starts the scraping process with the given concurrency and timebetween
func startScraping(db *database.Queries, concurrency int, timebetween time.Duration) {
	log.Printf("Scraping on %v threads every %s", concurrency, timebetween)
	// Create a ticker that ticks every timebetween
	ticker := time.NewTicker(timebetween)
	// Loop that runs every time the ticker ticks
	for ; ; <-ticker.C {
		// Get the next feeds to fetch
		feeds, err := db.GetNextFeedsToFetch(context.Background(), int32(concurrency))
		if err != nil {
			log.Printf("Error fetching feeds: %v", err)
			continue
		}
		wg := &sync.WaitGroup{}
		for _, feed := range feeds {
			// Increment the wait group
			wg.Add(1)

			go scrapeFeed(db, feed, wg)
		}
		wg.Wait()
	}
}

func scrapeFeed(db *database.Queries, feed database.Feed, wg *sync.WaitGroup) {
	defer wg.Done()

	// Mark the feed as fetched
	_, err := db.MarkFeedAsFetched(context.Background(), feed.ID)
	if err != nil {
		log.Printf("Error marking feed as fetched: %v", err)
		return
	}
	// Scrape the feed
	rssFeed, err := urlToFeed(feed.Url)
	if err != nil {
		log.Printf("Error fetching feed: %v", err)
		return
	}

	for _, item := range rssFeed.Channel.Items {
		description := sql.NullString{}
		if item.Description != "" {
			description = sql.NullString{
				String: item.Description,
				Valid:  true,
			}
		}

		publishedAt, err := time.Parse(time.RFC1123Z, item.PubDate)
		if err != nil {
			log.Printf("Error parsing date: %v", err)
			continue
		}

		_, err = db.CreatePost(context.Background(), database.CreatePostParams{
			ID:          uuid.New(),
			CreatedAt:   time.Now(),
			UpdatedAt:   time.Now(),
			Title:       item.Title,
			Description: description,
			PublishedAt: publishedAt,
			Url:         item.Link,
			FeedID:      feed.ID,
		})
		if err != nil {
			if strings.Contains(err.Error(), `duplicate key value violates unique constraint "posts_url_key"`) {
				continue
			}
			log.Printf("Error creating post: %v", err)
		}
	}

	log.Printf("Feed %v collected, %v posts found", feed.Url, len(rssFeed.Channel.Items))
}
