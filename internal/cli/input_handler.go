package cli

import (
	"fmt"
	"strings"
)

func Execute(state *State, args []string) error {
	if len(args) < 2 {
		return ErrNoCommand
	}

	cmdName := strings.ToLower(args[1]) // Parseo a minusculas para evitar problemas de case-sensitive
	cmdArgs := args[2:]                 // Los argumentos del comando son todos los elementos después del nombre del comando

	cmd := command{
		Name: cmdName,
		args: cmdArgs,
	}

	commands := NewCommands()
	commands.register("login", handlerLogin)
	commands.register("register", handlerRegister)
	commands.register("reset", handlerReset)
	commands.register("users", handlerUsers)
	commands.register("agg", handlerAggregate)

	err := commands.run(state, cmd)
	if err != nil {
		return fmt.Errorf("error running command: %v", err)
	}
	return nil
}
