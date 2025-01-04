package endpoints

import (
	"context"

	"github.com/baothaihcmut/Ecommerce-Go/users/internal/core/commands"
	"github.com/baothaihcmut/Ecommerce-Go/users/internal/core/port/inbound"
	"github.com/go-kit/kit/endpoint"
)

type UserEnpoints struct {
	CreateUser   endpoint.Endpoint
	FindUserById endpoint.Endpoint
}

func MakeUserEndpoints(s inbound.UserService) UserEnpoints {
	return UserEnpoints{
		CreateUser: makeCreateUserEndpoint(s),
	}
}

func makeCreateUserEndpoint(s inbound.UserService) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		req := request.(*commands.CreateUserCommand)
		res, err := s.CreateUser(ctx, req)
		if err != nil {
			return nil, err
		}
		return res, nil
	}
}
