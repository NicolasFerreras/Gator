package cli

import (
	"context"
	"fmt"

	"github.com/NicolasFerreras/Gator/internal/database"
)

func handlerFollowin(state *State, cmd Command, user database.User) error {

	following, err := state.Db.GetFeedFollowsForUser(context.Background(), user.ID)
	if err != nil {
		return ErrDBGetFeeds(err)
	}

	for _, follow := range following {
		fmt.Println(follow.FeedName)
	}

	return nil
}
