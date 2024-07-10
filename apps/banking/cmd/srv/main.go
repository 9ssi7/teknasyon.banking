package main

import (
	"github.com/9ssi7/banking/api/rest"
	"github.com/9ssi7/banking/config"
	"github.com/9ssi7/banking/internal/app"
	"github.com/9ssi7/banking/internal/app/commands"
	"github.com/9ssi7/banking/internal/app/queries"
	"github.com/9ssi7/banking/internal/app/services"
	"github.com/9ssi7/banking/internal/domain/abstracts"
	"github.com/9ssi7/banking/internal/infra/db"
	"github.com/9ssi7/banking/internal/infra/db/migrations"
	"github.com/9ssi7/banking/internal/infra/db/seeds"
	"github.com/9ssi7/banking/internal/infra/keyval"
	"github.com/9ssi7/banking/internal/infra/repos"
	"github.com/9ssi7/banking/pkg/token"
	"github.com/9ssi7/banking/pkg/validation"
)

func main() {
	cnf := config.ReadValue()
	token.Init()
	db := db.ConnectDB()
	kvdb := keyval.ConnectDB()

	if cnf.Database.Migrate {
		migrations.Run(db)
	}
	if cnf.Database.Seed {
		seeds.Run(db)
	}

	r := abstracts.Repositories{
		UserRepo:        repos.NewUserRepo(db),
		AccountRepo:     repos.NewAccountRepo(db),
		TransactionRepo: repos.NewTransactionRepo(db),
		SessionRepo:     repos.NewSessionRepo(kvdb),
		VerifyRepo:      repos.NewVerifyRepo(kvdb),
	}

	v := validation.New()

	rest.New(app.App{
		Commands: commands.NewHandler(r, v),
		Queries:  queries.NewHandler(r),
		Services: services.NewHandler(),
	}).Listen()
}
