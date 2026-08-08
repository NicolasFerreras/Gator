package cli

import (
	"context"
	"fmt"
)

func handlerLogin(state *State, cmd command) error {
	if len(cmd.args) == 0 {
		return fmt.Errorf("no username provided")
	}

	username := cmd.args[0]
	// Check if the username exists in the database
	user, err := state.Db.GetUserByUsername(context.Background(), username)
	if err != nil {
		return fmt.Errorf("failed to get user: %v", err)
	}

	if user.Username == "" {
		return fmt.Errorf("username does not exist")
	}

	err = state.Config.SetUser(username)
	if err != nil {
		return fmt.Errorf("failed to set user: %v", err)
	}

	fmt.Printf("Successfully updated username to: %s\n", username)
	return nil
}
