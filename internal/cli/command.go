package cli

import "fmt"

type command struct {
	Name string
	args []string
}

func NewCommands() *Commands {
	return &Commands{
		cmdMap: make(map[string]func(*State, command) error),
	}
}

type Commands struct { // Handler para los comandos
	cmdMap map[string]func(*State, command) error
}

func (c *Commands) run(s *State, cmd command) error { // Ejecuta el comando correspondiente según el nombre del comando
	if handler, exists := c.cmdMap[cmd.Name]; exists {
		return handler(s, cmd)
	}
	return fmt.Errorf("unknown command: %s", cmd.Name)
}

func (c *Commands) register(name string, f func(*State, command) error) { // Registra un nuevo comando en el mapa de comandos
	if c.cmdMap == nil {
		c.cmdMap = make(map[string]func(*State, command) error)
	}
	c.cmdMap[name] = f

}
