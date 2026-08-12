package cli

import (
	"context"

	"fmt"
	"time"

	"github.com/NicolasFerreras/Gator/internal/database"
	"github.com/NicolasFerreras/Gator/internal/errors"
	"github.com/google/uuid"
)

func handlerFollow(state *State, cmd command) error {
	if len(cmd.args) < 1 {
		return ErrMissingFeedURL
	}

	feedURL := cmd.args[0]

	feed, err := state.Db.GetFeedByURL(context.Background(), feedURL)
	if err != nil {
		return errors.ErrDBGetFeeds(err)
	}

	currentUser, err := state.Db.GetUserID(context.Background(), state.Config.CurrentUserName)
	if err != nil {
		return err
	}

	_, err = state.Db.CreateFeedFollow(context.Background(), database.CreateFeedFollowParams{
		ID:        uuid.New(),
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		UserID:    currentUser,
		FeedID:    feed.ID,
	})
	if err != nil {
		return errors.ErrDBCreateFeedFollow(err)
	}

	m := fmt.Sprintf(
		`
Nombre: %s
Usuario: %s
			`, feed.Name, currentUser,
	)
	fmt.Println(m)

	return nil
}
