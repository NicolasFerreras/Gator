package cli

// Este archivo mantiene compatibilidad de imports antiguos.
// Los errores reales están en internal/errors/errors.go

import "github.com/NicolasFerreras/Gator/internal/errors_handling"

// Reexportar para uso directo: cli.ErrNoUsername, etc.
var (
	ErrNoUsername       = errors_handling.ErrNoUsername
	ErrNoCommand        = errors_handling.ErrNoCommand
	ErrUserExists       = errors_handling.ErrUserExists
	ErrUnknownCommand   = errors_handling.ErrUnknownCommand
	ErrNoArgument       = errors_handling.ErrNoArgument
	ErrTooManyArguments = errors_handling.ErrTooManyArguments

	ErrDBCheckUser        = errors_handling.ErrDBCheckUser
	ErrDBCreateUser       = errors_handling.ErrDBCreateUser
	ErrDBGetUser          = errors_handling.ErrDBGetUser
	ErrDBGetUsers         = errors_handling.ErrDBGetUsers
	ErrDBCreateFeed       = errors_handling.ErrDBCreateFeed
	ErrDBGetFeeds         = errors_handling.ErrDBGetFeeds
	ErrDBCreateFeedFollow = errors_handling.ErrDBCreateFeedFollow

	ErrConfigSetUser   = errors_handling.ErrConfigSetUser
	ErrNotLoggedIn     = errors_handling.ErrNotLoggedIn
	ErrMissingUserName = errors_handling.ErrMissingUserName
	ErrInvalidLimit    = errors_handling.ErrInvalidLimit

	ErrMakeRequest        = errors_handling.ErrMakeRequest
	ErrServer             = errors_handling.ErrServer
	ErrReadingResponse    = errors_handling.ErrReadingResponse
	ErrUnmarshal          = errors_handling.ErrUnmarshal
	ErrFetchFeed          = errors_handling.ErrFetchFeed
	ErrMissingFeedURL     = errors_handling.ErrMissingFeedURL
	ErrInvalidFeedURL     = errors_handling.ErrInvalidFeedURL
	ErrDBDeleteFeedFollow = errors_handling.ErrDBDeleteFeedFollow
)
