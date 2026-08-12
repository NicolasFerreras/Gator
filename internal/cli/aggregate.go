package cli

import (
	"context"
	"fmt"

	"github.com/NicolasFerreras/Gator/internal/errors"
	rss "github.com/NicolasFerreras/Gator/internal/rss"
)

const (
	feedURL = "https://www.wagslane.dev/index.xml"
)

func handlerAggregate(state *State, cmd Command) error {
	// Check if the user is logged in
	if state.Config.CurrentUserName == "" {
		return errors.ErrNotLoggedIn
	}

	ctx := context.Background()

	// Fetch the RSS feed
	feed, err := rss.FetchFeed(ctx, feedURL)
	if err != nil {
		return errors.ErrFetchFeed(err)
	}

	// Print the fetched items
	for _, item := range feed.Channel.Item {
		fmt.Printf("Title: %s\n", item.Title)
		fmt.Printf("Link: %s\n", item.Link)
		fmt.Printf("Description: %s\n", item.Description)
		fmt.Println()
	}

	return nil
}
