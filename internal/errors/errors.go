package errors

import "fmt"

// Mensajes de error centralizados para TODO el proyecto
const (
	// Argumentos y validación de entrada
	noUsernameMsg  = "no username provided"
	noCommandMsg   = "no command provided"
	userExistsMsg  = "username already exists"
	unknownCommand = "unknown command: %s"

	// Errores de base de datos (formato)
	fmtDBCheckUser  = "failed to check existing user: %w"
	fmtDBCreateUser = "failed to create user: %w"
	fmtDBGetUser    = "failed to get user: %v"
	fmtDBGetUsers   = "failed to get users: %v"

	// Errores de configuración
	fmtConfigSetUser = "failed to set user: %v"
	fmtNotLoggedIn   = "no user is currently logged in"

	// Errores de RSS
	fmtMakeRequest        = "failed to make HTTP request: %v"
	fmtServer             = "server error: %v"
	fmtReadingResponse    = "failed to read response body: %v"
	fmtUnmarshal          = "failed to unmarshal XML: %v"
	fmtFetchFeed          = "failed to fetch feed: %v"
	fmtDBDeleteFeedFollow = "failed to delete feed follow: %v"
)

// Variables de error predefinidas para uso directo
var (
	ErrNoUsername     = fmt.Errorf(noUsernameMsg)
	ErrNoCommand      = fmt.Errorf(noCommandMsg)
	ErrUserExists     = fmt.Errorf(userExistsMsg)
	ErrUnknownCommand = func(cmd string) error {
		return fmt.Errorf(unknownCommand, cmd)
	}

	// Factories para errores envueltos con contexto DB
	ErrDBCheckUser        = func(err error) error { return fmt.Errorf(fmtDBCheckUser, err) }
	ErrDBCreateUser       = func(err error) error { return fmt.Errorf(fmtDBCreateUser, err) }
	ErrDBGetUser          = func(err error) error { return fmt.Errorf(fmtDBGetUser, err) }
	ErrDBGetUsers         = func(err error) error { return fmt.Errorf(fmtDBGetUsers, err) }
	ErrDBCreateFeed       = func(err error) error { return fmt.Errorf("failed to create feed: %v", err) }
	ErrDBGetFeeds         = func(err error) error { return fmt.Errorf("failed to get feeds: %v", err) }
	ErrDBCreateFeedFollow = func(err error) error { return fmt.Errorf("failed to create feed follow: %v", err) }

	// Errores de configuración
	ErrConfigSetUser   = func(err error) error { return fmt.Errorf(fmtConfigSetUser, err) }
	ErrNotLoggedIn     = fmt.Errorf(fmtNotLoggedIn)
	ErrMissingUserName = fmt.Errorf("missing username argument")

	// Errores de RSS
	ErrMakeRequest     = func(err error) error { return fmt.Errorf(fmtMakeRequest, err) }
	ErrServer          = func(err error) error { return fmt.Errorf("server error: %v", err) }
	ErrReadingResponse = func(err error) error { return fmt.Errorf("failed to read response body: %v", err) }
	ErrUnmarshal       = func(err error) error { return fmt.Errorf("failed to unmarshal XML: %v", err) }
	ErrFetchFeed       = func(err error) error { return fmt.Errorf(fmtFetchFeed, err) }
	ErrMissingFeedURL  = fmt.Errorf("missing feed URL argument")
	ErrInvalidFeedURL  = func(err error) error {
		return fmt.Errorf("Invalid argument, expected name then URL, but got only one argument: %v", err)
	}
	ErrDBDeleteFeedFollow = func(err error) error { return fmt.Errorf("failed to delete feed follow: %v", err) }
)
