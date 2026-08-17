package cli

import (
	"context"
	"fmt"
	"strconv"

	"github.com/NicolasFerreras/Gator/internal/database"
	"github.com/NicolasFerreras/Gator/internal/errors_handling"
)

func handlerBrowse(state *State, cmd Command) error {
	var optionalLimit int // Default limit
	if len(cmd.args) == 0 {
		optionalLimit = 2
	}
	if len(cmd.args) == 1 {
		limit, err := strconv.Atoi(cmd.args[0])
		if err != nil {
			return errors_handling.ErrInvalidLimit(err)
		}
		optionalLimit = limit
	}
	if len(cmd.args) > 1 {
		return errors_handling.ErrTooManyArguments
	}

	userId, err := state.Db.GetUserID(context.Background(), state.Config.CurrentUserName)
	if err != nil {
		return err
	}
	post, err := state.Db.GetPostsByUserId(context.Background(), database.GetPostsByUserIdParams{
		UserID: userId,
		Limit:  int32(optionalLimit),
		Offset: 0,
	})
	if err != nil {
		return err
	}
	for _, p := range post {
		m := `
Nombre: %s
Descripcion: %s
URL: %s
Fecha de publicacion: %s
		`
		m = m + "\n"
		m = m + "\n"
		m = m + "----------------------------------------"

		fmt.Printf(m, p.Title, p.Description, p.Url, p.PublishedAt)
	}

	return nil
}
