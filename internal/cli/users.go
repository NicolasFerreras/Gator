package cli

import (
	"context"
	"fmt"
)

func handlerUsers(state *State, cmd command) error {
	users, err := state.Db.GetUsers(context.Background())
	if err != nil {
		return ErrDBGetUsers(err)
	}
	for _, user := range users {
		if user.Username == state.Config.CurrentUserName {
			fmt.Printf("* %s (current)\n", user.Username)
		} else {
			fmt.Printf("* %s\n", user.Username)
		}
	}
	return nil
}
