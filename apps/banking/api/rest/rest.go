package rest

import (
	"context"
	"fmt"

	"github.com/9ssi7/banking/api/rest/middlewares"
	"github.com/9ssi7/banking/api/rest/routes"
	restsrv "github.com/9ssi7/banking/api/rest/srv"
	"github.com/9ssi7/banking/config"
	"github.com/9ssi7/banking/internal/app"
	"github.com/9ssi7/banking/pkg/server"
	"github.com/goccy/go-json"
	"github.com/gofiber/contrib/otelfiber/v2"
	"github.com/gofiber/fiber/v2"
	"go.opentelemetry.io/otel/metric"
	"go.opentelemetry.io/otel/trace"
)

type srv struct {
	app    app.App
	fiber  *fiber.App
	srv    restsrv.Srv
	meter  metric.Meter
	tracer trace.Tracer
}

type Config struct {
	App    app.App
	Meter  metric.Meter
	Tracer trace.Tracer
}

func New(cnf Config) server.Listener {
	restsrv := restsrv.New(cnf.App)
	return srv{
		app:    cnf.App,
		meter:  cnf.Meter,
		tracer: cnf.Tracer,
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
	requestDuration, err := s.meter.Float64Histogram(
		"http_request_duration",
		metric.WithDescription("HTTP req duration"),
		metric.WithUnit("ms"),
	)
	if err != nil {
		return err
	}
	requestCount, err := s.meter.Int64Counter(
		"http_request_count",
		metric.WithDescription("HTTP req count"),
	)
	if err != nil {
		return err
	}
	s.fiber.Use(s.srv.Cors(), s.srv.IpAddr())
	s.fiber.Use(otelfiber.Middleware(otelfiber.WithServerName("banking"), otelfiber.WithCollectClientIP(true)))
	s.fiber.Use(middlewares.Metric(requestDuration, requestCount, s.tracer))
	s.fiber.Use(s.srv.DeviceId())
	routes.Auth(s.fiber, s.srv, s.app)
	routes.Account(s.fiber, s.srv, s.app)
	routes.Transaction(s.fiber, s.srv, s.app)
	routes.MoneyTransfer(s.fiber, s.srv, s.app)
	return s.fiber.Listen(fmt.Sprintf("%v:%v", configs.Http.Host, configs.Http.Port))
}

func (s srv) Shutdown(ctx context.Context) error {
	return s.fiber.ShutdownWithContext(ctx)
}
