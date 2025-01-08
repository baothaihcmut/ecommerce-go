package endpoints

import (
	"context"

	"github.com/baothaihcmut/Ecommerce-Go/users/internal/core/port/inbound/command/commands"
	commandHandler "github.com/baothaihcmut/Ecommerce-Go/users/internal/core/port/inbound/command/handlers"
	"github.com/go-kit/kit/endpoint"
)

type AuthEndpoints struct {
	Login endpoint.Endpoint
}

func MakeAuthEndpoints(c commandHandler.AuthCommandHandler) AuthEndpoints {
	return AuthEndpoints{
		Login: makeLoginEndpoint(c),
	}
}

func makeLoginEndpoint(c commandHandler.AuthCommandHandler) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		req := request.(*commands.LoginCommand)
		res, err := c.Login(ctx, req)
		if err != nil {
			return nil, err
		}
		return res, nil
	}
}
