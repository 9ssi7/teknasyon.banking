package observer

import (
	"context"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/exporters/otlp/otlpmetric/otlpmetricgrpc"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracegrpc"
	"go.opentelemetry.io/otel/metric"
	"go.opentelemetry.io/otel/propagation"
	sdkmetric "go.opentelemetry.io/otel/sdk/metric"
	"go.opentelemetry.io/otel/sdk/resource"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.12.0"
	"go.opentelemetry.io/otel/trace"
	"google.golang.org/grpc/credentials"
)

type Service interface {
	Shutdown(ctx context.Context) error
	Init(ctx context.Context) error
	GetMeter() metric.Meter
	GetTracer() trace.Tracer
}

type Config struct {
	Name     string
	Endpoint string
	UseSSL   bool
}

type srv struct {
	cfg           Config
	res           *resource.Resource
	traceProvider *sdktrace.TracerProvider
	meterProvider *sdkmetric.MeterProvider

	meter  metric.Meter
	tracer trace.Tracer
}

func New(cfg Config) Service {
	return &srv{
		cfg: cfg,
	}
}

func (s *srv) initResource(ctx context.Context) error {
	res, err := resource.New(ctx,
		resource.WithAttributes(
			semconv.ServiceNameKey.String(s.cfg.Name),
		),
	)
	if err != nil {
		return err
	}
	s.res = res
	return nil
}

func (s *srv) initTracerProvider(ctx context.Context) error {
	options := []otlptracegrpc.Option{otlptracegrpc.WithEndpoint(s.cfg.Endpoint)}
	if s.cfg.UseSSL {
		options = append(options, otlptracegrpc.WithTLSCredentials(credentials.NewClientTLSFromCert(nil, "")))
	} else {
		options = append(options, otlptracegrpc.WithInsecure())
	}
	traceExporter, err := otlptrace.New(ctx, otlptracegrpc.NewClient(options...))
	if err != nil {
		return err
	}
	bsp := sdktrace.NewBatchSpanProcessor(traceExporter)
	tracerProvider := sdktrace.NewTracerProvider(
		sdktrace.WithSampler(sdktrace.AlwaysSample()),
		sdktrace.WithResource(s.res),
		sdktrace.WithSpanProcessor(bsp),
	)
	otel.SetTracerProvider(tracerProvider)
	otel.SetTextMapPropagator(propagation.NewCompositeTextMapPropagator(
		propagation.TraceContext{},
		propagation.Baggage{},
	))
	s.traceProvider = tracerProvider
	return nil
}

func (s *srv) initMeterProvider(ctx context.Context) error {
	options := []otlpmetricgrpc.Option{otlpmetricgrpc.WithEndpoint(s.cfg.Endpoint)}
	if s.cfg.UseSSL {
		options = append(options, otlpmetricgrpc.WithTLSCredentials(credentials.NewClientTLSFromCert(nil, "")))
	} else {
		options = append(options, otlpmetricgrpc.WithInsecure())
	}
	metricExporter, err := otlpmetricgrpc.New(ctx, options...)
	if err != nil {
		return err
	}
	meterProvider := sdkmetric.NewMeterProvider(
		sdkmetric.WithResource(s.res),
		sdkmetric.WithReader(sdkmetric.NewPeriodicReader(metricExporter)),
	)
	otel.SetMeterProvider(meterProvider)
	s.meterProvider = meterProvider
	return nil
}

func (s *srv) Init(ctx context.Context) error {
	if err := s.initResource(ctx); err != nil {
		return err
	}
	if err := s.initTracerProvider(ctx); err != nil {
		return err
	}
	if err := s.initMeterProvider(ctx); err != nil {
		return err
	}
	s.tracer = otel.Tracer(s.cfg.Name)
	s.meter = otel.Meter(s.cfg.Name)
	return nil
}

func (s *srv) Shutdown(ctx context.Context) error {
	if err := s.traceProvider.Shutdown(ctx); err != nil {
		return err
	}
	if err := s.meterProvider.Shutdown(ctx); err != nil {
		return err
	}
	return nil
}

func (s *srv) GetMeter() metric.Meter {
	return s.meter
}

func (s *srv) GetTracer() trace.Tracer {
	return s.tracer
}
