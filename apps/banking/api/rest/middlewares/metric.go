package middlewares

import (
	"time"

	"github.com/gofiber/fiber/v2"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/metric"
	"go.opentelemetry.io/otel/trace"
)

func Metric(durationM metric.Float64Histogram, reqM metric.Int64Counter, tracer trace.Tracer) fiber.Handler {
	return func(c *fiber.Ctx) error {
		start := time.Now()
		err := c.Next()
		elapsed := time.Since(start).Seconds()
		durationM.Record(c.UserContext(), elapsed, metric.WithAttributes(attribute.String("method", c.Method()), attribute.String("path", c.Path())))
		reqM.Add(c.UserContext(), 1, metric.WithAttributes(attribute.String("method", c.Method()), attribute.String("path", c.Path())))
		return err
	}
}
