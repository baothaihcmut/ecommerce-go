package transports

import (
	"context"

	"github.com/baothaihcmut/Ecommerce-Go/users/internal/adapter/grpc/proto"
	"github.com/go-kit/log"
	"google.golang.org/grpc/codes"
)

func ErrorHandle(logger log.Logger) func(ctx context.Context, req interface{}, resp GrpcResponse, err error) (interface{}, error) {
	return func(ctx context.Context, req interface{}, resp GrpcResponse, err error) (interface{}, error) {
		if err != nil {
			logger.Log("event", "error", "error", err)
			errorResp := &proto.Status{
				Success: false,
				Code:    mapErrorToGrpcStatus(err).String(),
				Message: err.Error(),
				Details: []string{err.Error()},
			}
			resp.Status = errorResp
			return resp, nil
		}
		resp.Status = &proto.Status{
			Success: true,
			Code:    codes.OK.String(),
			Message: "Success",
		}
		return resp, nil
	}
}
