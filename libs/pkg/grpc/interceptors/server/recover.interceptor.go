package interceptors

import (
	"context"
	"errors"
	"runtime"

	"github.com/baothaihcmut/Ecommerce-Go/libs/pkg/tracing"
	"go.opentelemetry.io/otel/trace"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func RecoverInterceptor(tracer trace.Tracer) func(ctx context.Context,
	req interface{},
	info *grpc.UnaryServerInfo,
	handler grpc.UnaryHandler) (interface{}, error) {
	return func(ctx context.Context,
		req interface{},
		info *grpc.UnaryServerInfo,
		handler grpc.UnaryHandler) (_ interface{}, err error) {
		ctx, span := tracing.StartSpan(ctx, tracer, "Grpc server: interceptor", nil)
		defer tracing.EndSpan(span, err, nil)
		defer func() {
			if r := recover(); r != nil {
				//convert error to string
				var detail string
				switch r.(type) {
				case string:
					detail = r.(string)
				case error:
					detail = r.(error).Error()
				default:
					detail = "Cannot convert error"
				}
				buf := make([]byte, 1024)
				n := runtime.Stack(buf, false)
				// tracing
				tracing.SetSpanAttribute(span, map[string]interface{}{
					"detail":     detail,
					"level":      "FATAL",
					"stackTrace": string(buf[:n]),
				})
				tracing.EndSpan(span, errors.New("Panic error"), nil)
				err = status.Error(codes.Internal, "Internal error")
			}
		}()
		res, err := handler(ctx, req)
		return res, err
	}
}
