package interceptors

import (
	"context"

	"github.com/baothaihcmut/Ecommerce-go/libs/pkg/constant"
	"github.com/baothaihcmut/Ecommerce-go/libs/pkg/logger"
	"github.com/google/uuid"
	"google.golang.org/grpc"
)



func LoggingInterceptor(logger logger.Logger) func (
	ctx context.Context,
	req any,
	info *grpc.UnaryServerInfo,
	handler grpc.UnaryHandler,
) (any, error)  {
	return func(
		ctx context.Context, 
		req any, 
		info *grpc.UnaryServerInfo, 
		handler grpc.UnaryHandler) (any, error) {
		ctx = context.WithValue(ctx, constant.RequestIdKey, uuid.New())
		logger.WithCtx(ctx).Info(map[string]any{
			"method": info.FullMethod,
		},"Incomming request")
		res,err := handler(ctx,req)
		if err != nil {
			logger.WithCtx(ctx).Errorf(map[string]any{
				"detail": err.Error(),
			},"Error")
			return nil,err
		}
		logger.WithCtx(ctx).Info(map[string]any{
			"method": info.FullMethod,

		},"Outgoing response")
		return res,nil
	}
}