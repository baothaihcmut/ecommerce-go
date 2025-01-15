package interceptor

import (
	"context"
	"fmt"

	"github.com/baothaihcmut/Ecommerce-Go/api-gateway/internal/common/grpc/proto"
	"github.com/baothaihcmut/Ecommerce-Go/api-gateway/internal/common/utils"
	"github.com/labstack/echo/v4"
	"github.com/mitchellh/mapstructure"
	"google.golang.org/grpc"
)

type GrpcResponse[T any] struct {
	Status *proto.Status `mapstruct:"Status"`
	Data   *T            `mapstructure:"Data"`
}

func ErrorHandlerClientInterceptor[T any]() grpc.UnaryClientInterceptor {
	return func(
		ctx context.Context,
		method string,
		req, reply interface{},
		cc *grpc.ClientConn,
		invoker grpc.UnaryInvoker,
		opts ...grpc.CallOption,
	) error {
		// Call the next interceptor in the chain or the actual invoker
		err := invoker(ctx, method, req, reply, cc, opts...)
		if err != nil {
			return err
		}

		var resp GrpcResponse[T]
		if err := mapstructure.Decode(reply, &resp); err != nil {
			return fmt.Errorf("failed to decode response: %w", err)
		}

		// Perform custom logic or validation on the response
		if !resp.Status.Success {
			return echo.NewHTTPError(utils.MapGrpcCodeToHttpCode(resp.Status.Code), resp.Status.Message)
		}
		return nil
	}
}
