package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/signal"
	"time"

	"github.com/9ssi7/banking/api/rest"
	"github.com/9ssi7/banking/config"
	"github.com/9ssi7/banking/internal/app"
	"github.com/9ssi7/banking/internal/app/commands"
	"github.com/9ssi7/banking/internal/app/queries"
	"github.com/9ssi7/banking/internal/app/services"
	"github.com/9ssi7/banking/internal/domain/abstracts"

	"github.com/9ssi7/banking/internal/infra/db/migrations"
	"github.com/9ssi7/banking/internal/infra/db/seeds"
	"github.com/9ssi7/banking/internal/infra/repos"
	"github.com/9ssi7/banking/pkg/retry"
	"github.com/9ssi7/banking/pkg/token"
	"github.com/9ssi7/banking/pkg/validation"
	"github.com/redis/go-redis/v9"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func main() {
	config.ReadValue()
	config.ReadValue()
	config.ReadValue()
	cnf := config.ReadValue()
	token.Init()
	db, err := connectPostgres(cnf.Database)
	if err != nil {
		panic(err)
	}
	kvdb, err := connectKeyVal(cnf.KeyValueDb)
	if err != nil {
		panic(err)
	}
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

	restHttp := rest.New(app.App{
		Commands: commands.NewHandler(r, v),
		Queries:  queries.NewHandler(r, v),
		Services: services.NewHandler(),
	})

	shutdownCh := make(chan os.Signal, 1)
	signal.Notify(shutdownCh, os.Interrupt)

	go func() {
		<-shutdownCh
		log.Println("restHttp is shutting down...")
		timeout := 5 * time.Second
		if err := restHttp.Shutdown(timeout); err != nil {
			log.Printf("restHttp shutdown failed: %v", err)
		}
	}()

	if err := restHttp.Listen(); err != nil {
		log.Fatalf("restHttp Listen failed: %v", err)
	}
}

func connectPostgres(cfg config.Database) (*gorm.DB, error) {
	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%d sslmode=%s",
		cfg.Host, cfg.User, cfg.Password, cfg.Name, cfg.Port, cfg.SslMode)
	var db *gorm.DB
	var err error
	err = retry.Run(func() error {
		db, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})
		return err
	}, retry.DefaultConfig)
	if err != nil {
		return nil, err
	}
	return db, nil
}

func connectKeyVal(cfg config.KeyValueDb) (*redis.Client, error) {
	rClient := redis.NewClient(&redis.Options{
		Addr:     cfg.Host + ":" + cfg.Port,
		Password: cfg.Pw,
		DB:       cfg.Db,
	})
	err := retry.Run(func() error {
		return rClient.Ping(context.Background()).Err()
	}, retry.DefaultConfig)
	if err != nil {
		return nil, err
	}
	return rClient, nil
}
