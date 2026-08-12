package cli

import (
	"context"
	"fmt"
)

func handlerUsers(state *State, cmd Command) error {
	users, err := state.Db.GetUsers(context.Background())
	if err != nil {
		return ErrDBGetUsers(err)
	}
	currentUser, err := GetCurrentUser(state)
	if err != nil {
		return fmt.Errorf("failed to get current user: %v", err)
	}
	for _, user := range users {
		if user.Username == currentUser {
			fmt.Printf("* %s (current)\n", user.Username)
		} else {
			fmt.Printf("* %s\n", user.Username)
		}
	}
	return nil
}

func GetCurrentUser(state *State) (string, error) {
	currentUser := state.Config.CurrentUserName
	if currentUser == "" {
		return "", fmt.Errorf("no current user set")
	}
	currentUserData, err := state.Db.GetUserByUsername(context.Background(), currentUser)
	if err != nil {
		return "", fmt.Errorf("failed to get current user from database: %v", err)
	}
	return currentUserData.Username, nil
}
