package client

import (
	"context"

	"github.com/baothaihcmut/Ecommerce-Go/libs/pkg/constance"
	"github.com/baothaihcmut/Ecommerce-Go/libs/pkg/models"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
)

func InjectAuthInterceptor() grpc.UnaryClientInterceptor {
	return func(
		ctx context.Context,
		method string,
		req, reply any,
		cc *grpc.ClientConn,
		invoker grpc.UnaryInvoker,
		opts ...grpc.CallOption) error {
		// get user context
		userContext := ctx.Value(string(constance.UserContext)).(*models.UserContext)
		tokenContext := ctx.Value(string(constance.TokenContext)).(string)
		md := metadata.Pairs(
			"user-id", userContext.Id.String(),
			"user-role", string(userContext.Role),
			"user-token", tokenContext,
		)
		ctx = metadata.NewOutgoingContext(ctx, md)
		err := invoker(ctx, method, req, reply, cc, opts...)
		if err != nil {
			return err
		}
		return nil
	}
}
