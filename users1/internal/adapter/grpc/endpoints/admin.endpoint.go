package endpoints

import (
	"context"

	"github.com/baothaihcmut/Ecommerce-Go/libs/pkg/tracing"
	"github.com/baothaihcmut/Ecommerce-Go/users/internal/core/command/port/inbound/commands"
	"github.com/baothaihcmut/Ecommerce-Go/users/internal/core/command/port/inbound/handlers"
	"github.com/go-kit/kit/endpoint"
	"go.opentelemetry.io/otel/trace"
)

type AdminEndpoints struct {
	LogIn       endpoint.Endpoint
	VerifyToken endpoint.Endpoint
}

func makeAdminLogInEndpoint(adminCommandHandler handlers.AdminCommandHandler, tracer trace.Tracer) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		ctx, span := tracing.StartSpan(ctx, tracer, "Admin.LogIn: endpoint", nil)
		defer tracing.EndSpan(span, err, nil)
		req := request.(*commands.LoginCommand)
		res, err := adminCommandHandler.LogIn(ctx, req)
		if err != nil {
			return nil, err
		}
		return res, nil
	}
}

func makeAdminVerifyTokenEndpoint(adminCommandHandler handlers.AdminCommandHandler, tracer trace.Tracer) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		ctx, span := tracing.StartSpan(ctx, tracer, "Admin.VerifyToken: endpoint", nil)
		defer tracing.EndSpan(span, err, nil)
		req := request.(*commands.VerifyTokenCommand)
		res, err := adminCommandHandler.VerifyToken(ctx, req)
		if err != nil {
			return nil, err
		}
		return res, nil
	}
}

func MakeAdminEndpoints(adminCommandHandler handlers.AdminCommandHandler, tracer trace.Tracer) AdminEndpoints {
	return AdminEndpoints{
		LogIn:       makeAdminLogInEndpoint(adminCommandHandler, tracer),
		VerifyToken: makeAdminVerifyTokenEndpoint(adminCommandHandler, tracer),
	}
}
