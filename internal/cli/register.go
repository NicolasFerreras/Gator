package cli

import (
	"context"
	"fmt"
	"time"

	"database/sql"

	"github.com/NicolasFerreras/Gator/internal/database"
	"github.com/google/uuid"
)

func handlerRegister(state *State, cmd command) error {
	if len(cmd.args) == 0 {
		return fmt.Errorf("no username provided")
	}
	// Check if the username already exists in the database
	username := cmd.args[0]
	existingUser, err := state.Db.GetUserByUsername(context.Background(), username)
	if err != nil && err != sql.ErrNoRows {
		return fmt.Errorf("failed to check existing user: %v", err)
	}

	if err == nil && existingUser.Username != "" {
		return fmt.Errorf("username already exists")
	}
	if cmd.args[0] == existingUser.Username {
		return fmt.Errorf("username already exists")
	}

	// Create a new user in the database

	_, err = state.Db.CreateUser(context.Background(), database.CreateUserParams{
		ID:        uuid.New(),
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		Username:  username,
	})
	if err != nil {
		return fmt.Errorf("failed to create user: %v", err)
	}

	fmt.Printf("Successfully registered username: %s\n", cmd.args[0])

	// Update the current user in the config
	err = state.Config.SetUser(username)
	if err != nil {
		return fmt.Errorf("failed to set user: %v", err)
	}

	fmt.Printf("Successfully updated username to: %s\n", username)
	return nil
}
