package pkg

import (
	"context"
	"log/slog"
	"os"

	"go.opentelemetry.io/contrib/bridges/otelslog"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/exporters/otlp/otlpmetric/otlpmetricgrpc"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracehttp"
	"go.opentelemetry.io/otel/log/global"
	sdkmetric "go.opentelemetry.io/otel/sdk/metric"

	"go.opentelemetry.io/otel/exporters/otlp/otlplog/otlploghttp"
	sdklog "go.opentelemetry.io/otel/sdk/log"
	"go.opentelemetry.io/otel/sdk/resource"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.34.0"
)

type Otel struct{}

func NewOtel() *Otel {
	return &Otel{}
}

func (o *Otel) Init() {
	var (
		ctx         = context.Background()
		environment = os.Getenv("ENVIRONMENT")
	)

	if environment == "" {
		environment = "development"
	}

	res, err := resource.New(ctx, resource.WithAttributes(
		semconv.DeploymentEnvironmentName(environment),
		semconv.TelemetrySDKLanguageGo,
	))
	if err != nil {
		panic(err)
	}
	newTraceProvider(ctx, res)
	newLoggerProvider(ctx, res, environment)
	newMetricsProvider(ctx, res)
	slog.Info("Starting Otel")
}

func newTraceProvider(ctx context.Context, res *resource.Resource) {
	client := otlptracehttp.NewClient()
	exporter, err := otlptrace.New(ctx, client)
	if err != nil {
		panic(err)
	}
	tp := sdktrace.NewTracerProvider(
		sdktrace.WithBatcher(exporter),
		sdktrace.WithResource(res),
	)
	otel.SetTracerProvider(tp)
}

func newLoggerProvider(ctx context.Context, res *resource.Resource, environment string) {
	logExporter, err := otlploghttp.New(ctx)
	if err != nil {
		panic(err)
	}
	logProcessor := sdklog.NewBatchProcessor(logExporter)
	lp := sdklog.NewLoggerProvider(
		sdklog.WithProcessor(logProcessor),
		sdklog.WithResource(res),
	)

	var logger *slog.Logger

	if environment == "production" {
		logger = otelslog.NewLogger("bot-shopee-logger",
			otelslog.WithSource(true),
			otelslog.WithLoggerProvider(lp),
		)
	} else {
		jsonHandler := slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
			AddSource: true,
			Level:     slog.LevelInfo,
		})
		logger = slog.New(jsonHandler)
	}

	slog.SetDefault(logger)
	global.SetLoggerProvider(lp)
}

func newMetricsProvider(ctx context.Context, res *resource.Resource) {
	exporter, err := otlpmetricgrpc.New(ctx)
	if err != nil {
		panic(err)
	}
	mp := sdkmetric.NewMeterProvider(
		sdkmetric.WithReader(sdkmetric.NewPeriodicReader(exporter)),
		sdkmetric.WithResource(res),
	)
	otel.SetMeterProvider(mp)
}
