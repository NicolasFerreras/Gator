package cli

import (
	"github.com/NicolasFerreras/Gator/internal/config"
	"github.com/NicolasFerreras/Gator/internal/database"
)

type State struct {
	Config *config.Config
	Db     *database.Queries
}
