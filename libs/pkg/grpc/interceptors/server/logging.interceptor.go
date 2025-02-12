package interceptors

import (
	"context"

	"github.com/baothaihcmut/Ecommerce-Go/libs/pkg/logger"
	"google.golang.org/grpc"
)

func LoggingServerInterceptor(logger logger.ILogger) func(ctx context.Context,
	req interface{},
	info *grpc.UnaryServerInfo,
	handler grpc.UnaryHandler) (interface{}, error) {
	return func(ctx context.Context,
		req interface{},
		info *grpc.UnaryServerInfo,
		handler grpc.UnaryHandler) (interface{}, error) {
		resp, err := handler(ctx, req)
		if err != nil {
			return resp, nil
		}

		return resp, nil
	}
}
