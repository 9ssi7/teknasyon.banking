package rest

import (
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
	app app.App
}

func New(app app.App) server.Listener {
	return srv{app: app}
}

func (s srv) Listen() error {
	configs := config.ReadValue()
	restsrv := restsrv.New(s.app)
	app := fiber.New(fiber.Config{
		ErrorHandler:   restsrv.ErrorHandler(),
		AppName:        "banking",
		ServerHeader:   "banking",
		JSONEncoder:    json.Marshal,
		JSONDecoder:    json.Unmarshal,
		CaseSensitive:  true,
		BodyLimit:      10 * 1024 * 1024,
		ReadBufferSize: 10 * 1024 * 1024,
	})
	app.Use(restsrv.Cors(), restsrv.DeviceId(), restsrv.IpAddr())
	routes.Auth(app, restsrv, s.app)
	routes.Account(app, restsrv, s.app)
	routes.Transaction(app, restsrv, s.app)
	routes.MoneyTransfer(app, restsrv, s.app)
	return app.Listen(fmt.Sprintf("%v:%v", configs.Http.Host, configs.Http.Port))
}
