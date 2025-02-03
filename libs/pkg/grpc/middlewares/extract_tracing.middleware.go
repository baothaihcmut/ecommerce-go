package middlewares

import (
	"context"

	"github.com/go-kit/kit/endpoint"
	"github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/ext"
	"google.golang.org/grpc/metadata"
)

func ExtractTracingMiddleware(tracer opentracing.Tracer, operationName string) endpoint.Middleware {
	return func(next endpoint.Endpoint) endpoint.Endpoint {
		return func(ctx context.Context, request interface{}) (response interface{}, err error) {

			md, ok := metadata.FromIncomingContext(ctx)
			if !ok {
				return next(opentracing.ContextWithSpan(ctx, tracer.StartSpan(operationName)), request)
			}
			carrier := opentracing.TextMapCarrier{}
			for key, values := range md {
				if len(values) > 0 {
					carrier[key] = values[0]
				}
			}
			spanContext, err := tracer.Extract(opentracing.TextMap, carrier)
			if err != nil {
				return next(opentracing.ContextWithSpan(ctx, tracer.StartSpan(operationName)), request)
			}
			span := tracer.StartSpan(operationName, ext.RPCServerOption(spanContext))
			return next(opentracing.ContextWithSpan(ctx, span), request)
		}
	}
}
