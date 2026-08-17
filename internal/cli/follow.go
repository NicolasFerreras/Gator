package cli

import (
	"context"

	"fmt"
	"time"

	"github.com/NicolasFerreras/Gator/internal/database"
	"github.com/NicolasFerreras/Gator/internal/errors_handling"
	"github.com/google/uuid"
)

func handlerFollow(state *State, cmd Command, user database.User) error {
	if len(cmd.args) < 1 {
		return ErrMissingFeedURL
	}

	feedURL := cmd.args[0]

	feed, err := state.Db.GetFeedByURL(context.Background(), feedURL)
	if err != nil {
		return errors_handling.ErrDBGetFeeds(err)
	}

	currentUser, err := state.Db.GetUserID(context.Background(), state.Config.CurrentUserName)
	if err != nil {
		return err
	}

	_, err = state.Db.CreateFeedFollow(context.Background(), database.CreateFeedFollowParams{
		ID:        uuid.New(),
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		UserID:    user.ID,
		FeedID:    feed.ID,
	})
	if err != nil {
		return errors_handling.ErrDBCreateFeedFollow(err)
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
