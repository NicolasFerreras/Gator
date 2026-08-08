package cli

import (
	"fmt"
)

func handlerLogin(state *State, cmd command) error {
	if len(cmd.args) == 0 {
		return fmt.Errorf("no username provided")
	}

	username := cmd.args[0]
	err := state.Config.SetUser(username)
	if err != nil {
		return fmt.Errorf("failed to set user: %v", err)
	}

	fmt.Printf("Successfully updated username to: %s\n", username)
	return nil
}
