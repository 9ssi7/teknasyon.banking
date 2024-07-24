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
	"github.com/redis/go-redis/v9"
	"github.com/uptrace/opentelemetry-go-extra/otelgorm"
	"go.opentelemetry.io/otel"

	"github.com/9ssi7/banking/internal/infra/db/migrations"
	"github.com/9ssi7/banking/internal/infra/db/seeds"
	"github.com/9ssi7/banking/internal/infra/observer"
	"github.com/9ssi7/banking/internal/infra/repos"
	"github.com/9ssi7/banking/pkg/retry"
	"github.com/9ssi7/banking/pkg/timeouter"
	"github.com/9ssi7/banking/pkg/token"
	"github.com/9ssi7/banking/pkg/validation"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var (
	pgdb     *gorm.DB
	rdclient *redis.Client
	reps     abstracts.Repositories
	valSrv   validation.Service
	obsrvr   observer.Service
)

func init() {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	cnf := config.ReadValue()
	token.Init()
	obsrvr = observer.New(observer.Config{
		Name:     cnf.Observer.Name,
		Endpoint: cnf.Observer.Endpoint,
		UseSSL:   cnf.Observer.UseSSL,
	})
	if err := obsrvr.Init(ctx); err != nil {
		panic(err)
	}
	if err := connectPostgres(ctx, cnf.Database); err != nil {
		panic(err)
	}
	plugin := otelgorm.NewPlugin(otelgorm.WithTracerProvider(otel.GetTracerProvider()))
	if err := pgdb.Use(plugin); err != nil {
		panic(err)
	}
	if err := connectRedis(ctx, cnf.RedisDb); err != nil {
		panic(err)
	}
	if cnf.Database.Migrate {
		migrations.Run(pgdb)
	}
	if cnf.Database.Seed {
		seeds.Run(pgdb)
	}
	reps = abstracts.Repositories{
		UserRepo:        repos.NewUserRepo(pgdb),
		AccountRepo:     repos.NewAccountRepo(pgdb),
		TransactionRepo: repos.NewTransactionRepo(pgdb),
		SessionRepo:     repos.NewSessionRepo(rdclient),
		VerifyRepo:      repos.NewVerifyRepo(rdclient),
	}

	valSrv = validation.New()
}

func main() {
	tracer := obsrvr.GetTracer()
	a := app.App{
		Commands: commands.NewHandler(tracer, reps, valSrv),
		Queries:  queries.NewHandler(tracer, reps, valSrv),
		Services: services.NewHandler(),
	}
	restHttp := rest.New(rest.Config{
		App:    a,
		Meter:  obsrvr.GetMeter(),
		Tracer: tracer,
	})

	shutdownCh := make(chan os.Signal, 1)
	signal.Notify(shutdownCh, os.Interrupt)

	go func() {
		<-shutdownCh
		log.Println("application is shutting down...")
		timeout := 5 * time.Second
		ctx, cancel := context.WithTimeout(context.Background(), timeout)
		defer cancel()
		if err := restHttp.Shutdown(ctx); err != nil {
			log.Printf("restHttp shutdown failed: %v", err)
		}
		if err := disconnectPostgres(ctx); err != nil {
			log.Printf("pgdb close failed: %v", err)
		}
		if err := disconnectRedis(ctx); err != nil {
			log.Printf("rdb close failed: %v", err)
		}
		if err := obsrvr.Shutdown(ctx); err != nil {
			log.Printf("obsrvr shutdown failed: %v", err)
		}
	}()

	if err := restHttp.Listen(); err != nil {
		log.Fatalf("restHttp Listen failed: %v", err)
	}
}

func disconnectPostgres(ctx context.Context) error {
	return timeouter.Run(ctx, func() error {
		db, err := pgdb.DB()
		if err != nil {
			return err
		}
		return db.Close()
	})
}

func disconnectRedis(ctx context.Context) error {
	return timeouter.Run(ctx, func() error {
		return rdclient.Close()
	})
}

func connectPostgres(ctx context.Context, cfg config.Database) error {
	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%d sslmode=%s",
		cfg.Host, cfg.User, cfg.Password, cfg.Name, cfg.Port, cfg.SslMode)
	return retry.Run(func() error {
		return timeouter.Run(ctx, func() error {
			var err error
			pgdb, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})
			return err
		})
	}, retry.DefaultConfig)
}

func connectRedis(ctx context.Context, cfg config.RedisDb) error {
	rdclient = redis.NewClient(&redis.Options{
		Addr:     cfg.Host + ":" + cfg.Port,
		Password: cfg.Pw,
		DB:       cfg.Db,
	})
	return retry.Run(func() error {
		fmt.Println("connecting to redis...")
		return timeouter.Run(ctx, func() error {
			fmt.Println("pinging redis...")
			return rdclient.Ping(ctx).Err()
		})
	}, retry.DefaultConfig)
}
