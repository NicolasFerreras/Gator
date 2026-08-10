package cli

import (
	"context"
	"fmt"
)

func handlerLogin(state *State, cmd command) error {
	if len(cmd.args) == 0 {
		return ErrNoUsername
	}

	username := cmd.args[0]
	// Check if the username exists in the database
	user, err := state.Db.GetUserByUsername(context.Background(), username)
	if err != nil {
		return ErrDBGetUser(err)
	}

	if user.Username == "" {
		return fmt.Errorf("username does not exist")
	}

	err = state.Config.SetUser(username)
	if err != nil {
		return ErrConfigSetUser(err)
	}

	fmt.Printf("Logged in user %s\n", user.Username)
	return nil

}
