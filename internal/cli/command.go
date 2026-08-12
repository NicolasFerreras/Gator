package cli

type Command struct {
	Name string
	args []string
}

func NewCommands() *Commands {
	return &Commands{
		cmdMap: make(map[string]func(*State, Command) error),
	}
}

type Commands struct { // Handler para los comandos
	cmdMap map[string]func(*State, Command) error
}

func (c *Commands) run(s *State, cmd Command) error { // Ejecuta el comando correspondiente según el nombre del comando
	if handler, exists := c.cmdMap[cmd.Name]; exists {
		return handler(s, cmd)
	}
	return ErrUnknownCommand(cmd.Name)
}

func (c *Commands) register(name string, f func(*State, Command) error) { // Registra un nuevo comando en el mapa de comandos
	if c.cmdMap == nil {
		c.cmdMap = make(map[string]func(*State, Command) error)
	}
	c.cmdMap[name] = f

}
