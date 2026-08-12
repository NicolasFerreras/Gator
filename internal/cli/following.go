package cli

import (
	"context"
	"fmt"
)

func handlerFollowin(state *State, cmd command) error {
	currentUser, err := state.Db.GetUserID(context.Background(), state.Config.CurrentUserName)
	if err != nil {
		return err
	}
	following, err := state.Db.GetFeedFollowsForUser(context.Background(), currentUser)
	if err != nil {
		return ErrDBGetFeeds(err)
	}

	for _, follow := range following {
		fmt.Println(follow.FeedName)
	}

	return nil
}
