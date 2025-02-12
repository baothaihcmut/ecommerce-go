package interceptor

import (
	"context"
	"fmt"
	"net/http"

	"github.com/baothaihcmut/Ecommerce-Go/api-gateway/internal/common/grpc/proto"
	"github.com/baothaihcmut/Ecommerce-Go/api-gateway/internal/common/utils"
	"github.com/labstack/echo/v4"
	"github.com/mitchellh/mapstructure"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type GrpcResponse struct {
	Status *proto.Status `mapstruct:"Status"`
}

func ErrorHandlerClientInterceptor() grpc.UnaryClientInterceptor {
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
			//handle grpc error
			grpcStatus, ok := status.FromError(err)
			if ok {
				switch grpcStatus.Code() {
				case codes.Unavailable:
					return echo.NewHTTPError(http.StatusUnavailableForLegalReasons, "Service unavailable")
				default:
					return echo.NewHTTPError(http.StatusInternalServerError, "Internal server error")
				}
			}
			return err
		}

		var resp GrpcResponse
		if err := mapstructure.Decode(reply, &resp); err != nil {
			return fmt.Errorf("failed to decode response: %w", err)
		}

		// Perform custom logic or validation on the response
		if resp.Status.Code != codes.OK.String() {
			return echo.NewHTTPError(utils.MapGrpcCodeToHttpCode(resp.Status.Code), resp.Status.Message)
		}
		return nil
	}
}
