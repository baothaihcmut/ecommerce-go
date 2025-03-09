package endpoints

import (
	"context"

	"github.com/baothaihcmut/Ecommerce-Go/libs/pkg/tracing"
	"github.com/baothaihcmut/Ecommerce-Go/users/internal/core/command/port/inbound/commands"
	commandHandler "github.com/baothaihcmut/Ecommerce-Go/users/internal/core/command/port/inbound/handlers"
	"github.com/go-kit/kit/endpoint"
	"go.opentelemetry.io/otel/trace"
)

type AuthEndpoints struct {
	Login       endpoint.Endpoint
	SignUp      endpoint.Endpoint
	VerifyToken endpoint.Endpoint
}

func MakeAuthEndpoints(c commandHandler.AuthCommandHandler, tracer trace.Tracer) AuthEndpoints {
	return AuthEndpoints{
		Login:       makeLoginEndpoint(c, tracer),
		SignUp:      makeSignUpEndpoint(c, tracer),
		VerifyToken: makeVerifyTokenEndpoint(c, tracer),
	}
}

func makeLoginEndpoint(c commandHandler.AuthCommandHandler, tracer trace.Tracer) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		ctx, span := tracing.StartSpan(ctx, tracer, "Auth.LogIn: endpoint", nil)
		defer tracing.EndSpan(span, err, nil)
		req := request.(*commands.LoginCommand)
		res, err := c.Login(ctx, req)
		if err != nil {
			return nil, err
		}
		return res, nil
	}
}

func makeSignUpEndpoint(c commandHandler.AuthCommandHandler, tracer trace.Tracer) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		ctx, span := tracing.StartSpan(ctx, tracer, "Auth.SignUp: endpoint", nil)
		defer tracing.EndSpan(span, err, nil)
		req := request.(*commands.SignUpCommand)
		res, err := c.SignUp(ctx, req)
		if err != nil {
			return nil, err
		}
		return res, nil
	}
}
func makeVerifyTokenEndpoint(c commandHandler.AuthCommandHandler, tracer trace.Tracer) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		ctx, span := tracing.StartSpan(ctx, tracer, "Auth.VerifyToken: endpoint", nil)
		defer tracing.EndSpan(span, err, nil)
		req := request.(*commands.VerifyTokenCommand)
		res, err := c.VerifyToken(ctx, req)
		if err != nil {
			return nil, err
		}
		return res, nil
	}
}
