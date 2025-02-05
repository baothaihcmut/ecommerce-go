package tracing

import (
	"context"

	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/trace"
)

func setSpanAttributes(span trace.Span, attributes map[string]interface{}) {
	for key, val := range attributes {
		switch v := val.(type) {
		case string:
			span.SetAttributes(attribute.String(key, v))
		case []string:
			span.SetAttributes(attribute.StringSlice(key, v))
		case int:
			span.SetAttributes(attribute.Int(key, v))
		case []int:
			span.SetAttributes(attribute.IntSlice(key, v))
		case float64:
			span.SetAttributes(attribute.Float64(key, v))
		case []float64:
			span.SetAttributes(attribute.Float64Slice(key, v))
		case bool:
			span.SetAttributes(attribute.Bool(key, v))
		case []bool:
			span.SetAttributes(attribute.BoolSlice(key, v))
		default:
			// Handle unknown types
			// Optionally log or ignore unsupported attribute types
		}
	}
}

func StartSpan(ctx context.Context, tracer trace.Tracer, name string, attributes map[string]interface{}) (context.Context, trace.Span) {
	ctx, span := tracer.Start(ctx, name)
	setSpanAttributes(span, attributes)
	return ctx, span
}

func SetSpanAttribute(span trace.Span, attributes map[string]interface{}) {
	setSpanAttributes(span, attributes)
}

func EndSpan(span trace.Span, err error, attributes map[string]interface{}) {
	if err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())
	}
	setSpanAttributes(span, attributes)
	span.End()
}
