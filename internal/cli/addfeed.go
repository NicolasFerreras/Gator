package cli

import (
	"context"
	"net/url"

	"time"

	"github.com/NicolasFerreras/Gator/internal/database"
	"github.com/NicolasFerreras/Gator/internal/errors"
	"github.com/NicolasFerreras/Gator/internal/rss"
	"github.com/google/uuid"
)

func handlerAddFeed(state *State, cmd command) error {
	if len(cmd.args) < 2 {
		return ErrMissingFeedURL
	}

	_, err := url.ParseRequestURI(cmd.args[1])
	if err != nil {
		return errors.ErrInvalidFeedURL(err)
	}
	currentUser, err := state.Db.GetUserID(context.Background(), state.Config.CurrentUserName)
	if err != nil {
		return err
	}

	feedName := cmd.args[0]
	feedURL := cmd.args[1]
	feedUserId := currentUser

	ctx := context.Background()
	_, err = rss.FetchFeed(ctx, feedURL)
	if err != nil {
		return errors.ErrFetchFeed(err)
	}

	_, err = state.Db.CreateFeed(context.Background(), database.CreateFeedParams{
		ID:        uuid.New(),
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		Name:      feedName,
		Url:       feedURL,
		UserID:    feedUserId,
	})
	if err != nil {
		return errors.ErrDBCreateFeed(err)
	}
	return nil
}
