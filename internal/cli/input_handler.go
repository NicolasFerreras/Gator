package cli

import (
	"context"
	"fmt"
	"strings"

	"github.com/NicolasFerreras/Gator/internal/database"
)

func Execute(state *State, args []string) error {
	if len(args) < 2 {
		return ErrNoCommand
	}

	cmdName := strings.ToLower(args[1]) // Parseo a minusculas para evitar problemas de case-sensitive
	cmdArgs := args[2:]                 // Los argumentos del comando son todos los elementos después del nombre del comando

	cmd := Command{
		Name: cmdName,
		args: cmdArgs,
	}

	commands := NewCommands()
	commands.register("login", handlerLogin)
	commands.register("register", handlerRegister)
	commands.register("reset", handlerReset)
	commands.register("users", handlerUsers)
	commands.register("agg", handlerAggregate)
	commands.register("addfeed", middlewareLoggedIn(handlerAddFeed))
	commands.register("feeds", handlerFeeds)
	commands.register("follow", middlewareLoggedIn(handlerFollow))
	commands.register("following", middlewareLoggedIn(handlerFollowin))
	commands.register("unfollow", middlewareLoggedIn(handlerUnfollowFeed))

	err := commands.run(state, cmd)
	if err != nil {
		return fmt.Errorf("error running command: %v", err)
	}
	return nil
}

func middlewareLoggedIn(handler func(s *State, cmd Command, user database.User) error) func(*State, Command) error {
	return func(s *State, cmd Command) error {
		currentUser, err := s.Db.GetUserByUsername(context.Background(), s.Config.CurrentUserName)
		if err != nil {
			return err
		}
		return handler(s, cmd, currentUser)
	}
}
