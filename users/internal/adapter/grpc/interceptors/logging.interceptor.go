package interceptors

import (
	"context"

	"github.com/baothaihcmut/Ecommerce-Go/libs/pkg/logger"
	"google.golang.org/grpc"
)

func LoggingInterceptor(logger logger.ILogger) func(ctx context.Context,
	req interface{},
	info *grpc.UnaryServerInfo,
	handler grpc.UnaryHandler) (interface{}, error) {
	return func(ctx context.Context,
		req interface{},
		info *grpc.UnaryServerInfo,
		handler grpc.UnaryHandler) (interface{}, error) {
		resp, err := handler(ctx, req)
		if err != nil {
			logger.Error("err", err)
			return resp, nil
		}

		return resp, nil
	}
}
