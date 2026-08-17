package cli

import (
	"context"

	"github.com/NicolasFerreras/Gator/internal/database"
	"github.com/NicolasFerreras/Gator/internal/errors_handling"
)

func handlerUnfollowFeed(state *State, cmd Command, user database.User) error {
	if len(cmd.args) < 1 {
		return ErrMissingFeedURL
	}

	feedURL := cmd.args[0]

	feed, err := state.Db.GetFeedByURL(context.Background(), feedURL)
	if err != nil {
		return errors_handling.ErrDBGetFeeds(err)
	}

	err = state.Db.DeleteFeedFollowByUserID(context.Background(), database.DeleteFeedFollowByUserIDParams{
		UserID: user.ID,
		FeedID: feed.ID,
	})
	if err != nil {
		return errors_handling.ErrDBDeleteFeedFollow(err)
	}

	return nil
}
