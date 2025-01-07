package interceptors

import (
	"context"

	"github.com/go-kit/log"
	"github.com/go-kit/log/level"
	"google.golang.org/grpc"
)

func LoggingInterceptor(logger log.Logger) func(ctx context.Context,
	req interface{},
	info *grpc.UnaryServerInfo,
	handler grpc.UnaryHandler) (interface{}, error) {
	return func(ctx context.Context,
		req interface{},
		info *grpc.UnaryServerInfo,
		handler grpc.UnaryHandler) (interface{}, error) {
		level.Info(logger).Log("method", info.FullMethod, "request", req)
		resp, err := handler(ctx, req)
		if err != nil {
			level.Error(logger).Log("error", err)
			return resp, nil
		}

		return resp, nil
	}
}
