package cli

import (
	"context"
	"fmt"
)

func handlerFeeds(state *State, cmd command) error {
	ctx := context.Background()
	feeds, err := state.Db.GetFeedWithUserName(ctx)
	if err != nil {
		return ErrDBGetFeeds(err)
	}
	for _, feed := range feeds {
		m := fmt.Sprintf(
			`
Nombre: %s
URL: %s
Usuario: %s
			`, feed.Name, feed.Url, feed.Username,
		)
		fmt.Println(m)
	}
	return nil
}
