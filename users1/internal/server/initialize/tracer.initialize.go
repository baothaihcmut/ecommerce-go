package initialize

import (
	"context"
	"errors"
	"time"

	"github.com/baothaihcmut/Ecommerce-Go/users/internal/config"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracegrpc"
	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/sdk/resource"
	sdkTrace "go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.4.0"
	"go.opentelemetry.io/otel/trace"
)

func InitializeTracer(cfg *config.JaegerConfig) (*sdkTrace.TracerProvider, trace.Tracer, error) {
	exporter, err := otlptracegrpc.New(
		context.Background(),
		otlptracegrpc.WithEndpoint(cfg.Endpoint),
		otlptracegrpc.WithInsecure(),
	)
	if err != nil {
		return nil, nil, err
	}
	var sampler sdkTrace.Sampler
	if cfg.Sample == 100 {
		sampler = sdkTrace.AlwaysSample()
	} else if cfg.Sample < 100 && cfg.Sample > 0 {
		sampler = sdkTrace.TraceIDRatioBased(float64(cfg.Sample) / 100)
	} else if cfg.Sample == 0 {
		sampler = sdkTrace.NeverSample()
	} else {
		return nil, nil, errors.New("Trace sample must less than or equal 100 and greater or equal zero")
	}
	traceProvider := sdkTrace.NewTracerProvider(
		sdkTrace.WithSampler(sampler),
		sdkTrace.WithBatcher(exporter, sdkTrace.WithBatchTimeout(5*time.Second), sdkTrace.WithMaxExportBatchSize(50)),
		sdkTrace.WithResource(
			resource.NewWithAttributes(
				semconv.SchemaURL,
				semconv.ServiceNameKey.String("User service"),
			),
		),
	)
	otel.SetTracerProvider(traceProvider)
	otel.SetTextMapPropagator(propagation.TraceContext{})
	//get tracer
	return traceProvider, otel.Tracer("User Service"), nil
}
