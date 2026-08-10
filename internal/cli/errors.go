package cli

import "fmt"

// Mensajes de error centralizados para handlers CLI
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
	ErrDBCheckUser  = func(err error) error { return fmt.Errorf(fmtDBCheckUser, err) }
	ErrDBCreateUser = func(err error) error { return fmt.Errorf(fmtDBCreateUser, err) }
	ErrDBGetUser    = func(err error) error { return fmt.Errorf(fmtDBGetUser, err) }
	ErrDBGetUsers   = func(err error) error { return fmt.Errorf(fmtDBGetUsers, err) }

	// Errores de configuración
	ErrConfigSetUser = func(err error) error { return fmt.Errorf(fmtConfigSetUser, err) }
)
