package cli

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"

	"github.com/lib/pq"

	"github.com/NicolasFerreras/Gator/internal/database"

	"github.com/NicolasFerreras/Gator/internal/errors_handling"
	rss "github.com/NicolasFerreras/Gator/internal/rss"
	"github.com/google/uuid"
)

func scrapeFeeds(state *State) error {
	if state.Config.CurrentUserName == "" {
		return errors_handling.ErrNotLoggedIn
	}

	ctx := context.Background()
	nextFeed, err := state.Db.GetNextFeedToFetch(ctx)
	if err != nil {
		return err
	}

	feed, err := rss.FetchFeed(ctx, nextFeed.Url)
	if err != nil {
		return errors_handling.ErrFetchFeed(err)
	}

	state.Db.MarkFeedFetched(ctx, database.MarkFeedFetchedParams{
		UpdatedAt:     time.Now(),
		LastFetchedAt: time.Now(),
		ID:            nextFeed.ID,
	})
	for _, item := range feed.Channel.Item {
		var publishedAt sql.NullTime

		parsedTime, err := time.Parse(time.RFC1123Z, item.PubDate)
		publishedAt = sql.NullTime{Valid: true, Time: parsedTime}
		if err != nil {
			fmt.Printf("Error parsing time for item %s: %v\n", item.Title, err)
			publishedAt = sql.NullTime{Valid: false, Time: parsedTime}
		}

		_, err = state.Db.CreatePost(ctx, database.CreatePostParams{
			ID:          uuid.New(),
			CreatedAt:   time.Now(),
			UpdatedAt:   time.Now(),
			Title:       item.Title,
			Url:         item.Link,
			Description: sql.NullString{Valid: true, String: item.Description},
			PublishedAt: publishedAt,
			FeedID:      nextFeed.ID,
		})
		var pqErr *pq.Error
		if errors.As(err, &pqErr) && pqErr.Code == "23505" {
			continue
		}
		if err != nil {
			fmt.Printf("Error creating post: %v\n", err)
			continue
		}

	}
	return nil
}
