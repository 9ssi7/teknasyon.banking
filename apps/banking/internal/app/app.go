package app

import (
	"github.com/9ssi7/banking/internal/app/commands"
	"github.com/9ssi7/banking/internal/app/queries"
	"github.com/9ssi7/banking/internal/app/services"
)

type App struct {
	Commands commands.Handlers
	Queries  queries.Handlers
	Services services.Handlers
}
