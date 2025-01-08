package endpoints

import (
	"context"

	"github.com/baothaihcmut/Ecommerce-Go/users/internal/core/port/inbound/command/commands"
	commandHandler "github.com/baothaihcmut/Ecommerce-Go/users/internal/core/port/inbound/command/handlers"
	queryHandler "github.com/baothaihcmut/Ecommerce-Go/users/internal/core/port/inbound/query/handlers"
	"github.com/go-kit/kit/endpoint"
)

type UserEnpoints struct {
	CreateUser   endpoint.Endpoint
	FindUserById endpoint.Endpoint
}

func MakeUserEndpoints(c commandHandler.UserCommandHandler, q queryHandler.UserQueryHandler) UserEnpoints {
	return UserEnpoints{
		CreateUser: makeCreateUserEndpoint(c),
	}
}

func makeCreateUserEndpoint(s commandHandler.UserCommandHandler) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		req := request.(*commands.CreateUserCommand)
		res, err := s.CreateUser(ctx, req)
		if err != nil {
			return nil, err
		}
		return res, nil
	}
}
