package cli

import (
	"context"
	"fmt"
)

func handlerReset(state *State, cmd command) error {
	err := state.Db.DeleteAllUsers(context.Background())
	if err != nil {
		return fmt.Errorf("failed to delete all users: %v", err)
	}
	return nil
}
