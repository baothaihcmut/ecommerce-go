package endpoints

import (
	"context"

	"github.com/baothaihcmut/Ecommerce-Go/users/internal/core/commands"
	inboundCommand "github.com/baothaihcmut/Ecommerce-Go/users/internal/core/port/inbound/command"
	inboundQuery "github.com/baothaihcmut/Ecommerce-Go/users/internal/core/port/inbound/query"
	"github.com/go-kit/kit/endpoint"
)

type UserEnpoints struct {
	CreateUser   endpoint.Endpoint
	FindUserById endpoint.Endpoint
}

func MakeUserEndpoints(c inboundCommand.UserCommandPort, q inboundQuery.UserQueryPort) UserEnpoints {
	return UserEnpoints{
		CreateUser: makeCreateUserEndpoint(c),
	}
}

func makeCreateUserEndpoint(s inboundCommand.UserCommandPort) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		req := request.(*commands.CreateUserCommand)
		res, err := s.CreateUser(ctx, req)
		if err != nil {
			return nil, err
		}
		return res, nil
	}
}
