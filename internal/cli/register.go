package cli

import (
	"context"
	"fmt"
	"time"

	"database/sql"

	"github.com/NicolasFerreras/Gator/internal/database"
	"github.com/google/uuid"
)

func handlerRegister(state *State, cmd Command) error {
	if len(cmd.args) == 0 {
		return ErrNoUsername
	}

	username := cmd.args[0]

	// Check if the username already exists in the database
	existingUser, err := state.Db.GetUserByUsername(context.Background(), username)
	if err != nil && err != sql.ErrNoRows {
		return ErrDBCheckUser(err)
	}

	if err == nil && existingUser.Username != "" {
		return ErrUserExists
	}

	// Create a new user in the database
	_, err = state.Db.CreateUser(context.Background(), database.CreateUserParams{
		ID:        uuid.New(),
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		Username:  username,
	})
	if err != nil {
		return ErrDBCreateUser(err)
	}

	fmt.Printf("Successfully registered username: %s\n", username)

	// Update the current user in the config
	err = state.Config.SetUser(username)
	if err != nil {
		return ErrConfigSetUser(err)
	}

	fmt.Printf("Successfully updated username to: %s\n", username)
	return nil
}
