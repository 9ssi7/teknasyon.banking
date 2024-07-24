package rest

import (
	"context"
	"fmt"

	"github.com/9ssi7/banking/api/rest/routes"
	restsrv "github.com/9ssi7/banking/api/rest/srv"
	"github.com/9ssi7/banking/config"
	"github.com/9ssi7/banking/internal/app"
	"github.com/9ssi7/banking/pkg/server"
	"github.com/goccy/go-json"
	"github.com/gofiber/fiber/v2"
)

type srv struct {
	app   app.App
	fiber *fiber.App
	srv   restsrv.Srv
}

func New(app app.App) server.Listener {
	restsrv := restsrv.New(app)
	return srv{
		app: app,
		fiber: fiber.New(fiber.Config{
			ErrorHandler:   restsrv.ErrorHandler(),
			AppName:        "banking",
			ServerHeader:   "banking",
			JSONEncoder:    json.Marshal,
			JSONDecoder:    json.Unmarshal,
			CaseSensitive:  true,
			BodyLimit:      10 * 1024 * 1024,
			ReadBufferSize: 10 * 1024 * 1024,
		}),
		srv: restsrv,
	}
}

func (s srv) Listen() error {
	configs := config.ReadValue()

	s.fiber.Use(s.srv.Cors(), s.srv.DeviceId(), s.srv.IpAddr())
	routes.Auth(s.fiber, s.srv, s.app)
	routes.Account(s.fiber, s.srv, s.app)
	routes.Transaction(s.fiber, s.srv, s.app)
	routes.MoneyTransfer(s.fiber, s.srv, s.app)
	return s.fiber.Listen(fmt.Sprintf("%v:%v", configs.Http.Host, configs.Http.Port))
}

func (s srv) Shutdown(ctx context.Context) error {
	return s.fiber.ShutdownWithContext(ctx)
}
