package cli

// Este archivo mantiene compatibilidad de imports antiguos.
// Los errores reales están en internal/errors/errors.go

import "github.com/NicolasFerreras/Gator/internal/errors"

// Reexportar para uso directo: cli.ErrNoUsername, etc.
var (
	ErrNoUsername     = errors.ErrNoUsername
	ErrNoCommand      = errors.ErrNoCommand
	ErrUserExists     = errors.ErrUserExists
	ErrUnknownCommand = errors.ErrUnknownCommand

	ErrDBCheckUser       = errors.ErrDBCheckUser
	ErrDBCreateUser      = errors.ErrDBCreateUser
	ErrDBGetUser         = errors.ErrDBGetUser
	ErrDBGetUsers        = errors.ErrDBGetUsers

	ErrConfigSetUser     = errors.ErrConfigSetUser
	ErrDBGetLocationArea = errors.ErrDBGetLocationArea
	ErrNotLoggedIn       = errors.ErrNotLoggedIn

	ErrMakeRequest     = errors.ErrMakeRequest
	ErrServer          = errors.ErrServer
	ErrReadingResponse = errors.ErrReadingResponse
	ErrUnmarshal       = errors.ErrUnmarshal
	ErrFetchFeed       = errors.ErrFetchFeed
)
