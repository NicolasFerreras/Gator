package cli

import (
	"context"
	"fmt"
	"time"

	"github.com/NicolasFerreras/Gator/internal/database"
	"github.com/NicolasFerreras/Gator/internal/errors"
	rss "github.com/NicolasFerreras/Gator/internal/rss"
)

func handlerAggregate(state *State, cmd Command) error {
	if len(cmd.args) < 1 {
		return ErrNoArgument
	}

	reqTime, err := time.ParseDuration(cmd.args[0])
	if err != nil {
		return err
	}

	fmt.Printf("Collecting feeds every %v", reqTime)

	ticker := time.NewTicker(reqTime)
	//i := 0
	for ; ; <-ticker.C {
		scrapeFeeds(state)
		// i = i + 1
		// fmt.Printf("Fetch number %v", i) simple test para saber si esta haciendo el fetch correctamente

	}
}

func scrapeFeeds(state *State) error {
	if state.Config.CurrentUserName == "" {
		return errors.ErrNotLoggedIn
	}

	ctx := context.Background()
	nextFeed, err := state.Db.GetNextFeedToFetch(ctx)
	if err != nil {
		return err
	}

	feed, err := rss.FetchFeed(ctx, nextFeed.Url)
	if err != nil {
		return errors.ErrFetchFeed(err)
	}

	state.Db.MarkFeedFetched(ctx, database.MarkFeedFetchedParams{
		UpdatedAt:     time.Now(),
		LastFetchedAt: time.Now(),
		ID:            nextFeed.ID,
	})

	for _, item := range feed.Channel.Item {
		fmt.Printf("Title: %s\n", item.Title)
		fmt.Printf("Link: %s\n", item.Link)
		fmt.Printf("Description: %s\n", item.Description)
		fmt.Println()
	}

	return nil

}
