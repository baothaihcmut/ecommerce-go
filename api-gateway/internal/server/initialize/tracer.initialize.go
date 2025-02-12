package initialize

import (
	"context"
	"time"

	"github.com/baothaihcmut/Ecommerce-Go/api-gateway/internal/config"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracehttp"
	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/sdk/resource"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.4.0"
	"go.opentelemetry.io/otel/trace"
)

func InitializeTracer(cfg *config.Config) (*sdktrace.TracerProvider, trace.Tracer, error) {
	client := otlptracehttp.NewClient(otlptracehttp.WithEndpoint(cfg.JaegerConfig.Endpoint), otlptracehttp.WithInsecure())
	exporter, err := otlptrace.New(context.Background(), client)
	if err != nil {
		return nil, nil, err
	}
	tp := sdktrace.NewTracerProvider(
		sdktrace.WithSampler(sdktrace.AlwaysSample()),
		sdktrace.WithBatcher(exporter, sdktrace.WithBatchTimeout(5*time.Second), sdktrace.WithMaxExportBatchSize(50)),
		sdktrace.WithResource(
			resource.NewWithAttributes(
				semconv.SchemaURL,
				semconv.ServiceNameKey.String("Api gateway"),
			),
		),
	)
	otel.SetTracerProvider(tp)
	otel.SetTextMapPropagator(propagation.TraceContext{})
	return tp, otel.Tracer("Api gateway"), nil

}
