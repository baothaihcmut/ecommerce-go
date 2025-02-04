package middlewares

import (
	"context"

	"github.com/go-kit/kit/endpoint"
	"github.com/opentracing/opentracing-go"
)

func ExtractTracingMiddleware(tracer opentracing.Tracer, operationName string) endpoint.Middleware {
	return func(next endpoint.Endpoint) endpoint.Endpoint {
		return func(ctx context.Context, request interface{}) (response interface{}, err error) {
			span, ctx := opentracing.StartSpanFromContext(ctx, operationName)
			defer span.Finish()
			return next(opentracing.ContextWithSpan(ctx, span), request)
		}
	}
}
