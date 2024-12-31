package request

import (
	"context"

	"github.com/baothaihcmut/Ecommerce-Go/users/internal/adapter/grpc/proto"
	"github.com/baothaihcmut/Ecommerce-Go/users/internal/core/commands"
	"github.com/mitchellh/mapstructure"
)

type UserRequestMapper interface {
	ToCreateUserCommand(_ context.Context, request interface{}) (interface{}, error)
}

type UserRequestMapperImpl struct {
}

func (m *UserRequestMapperImpl) ToCreateUserCommand(_ context.Context, request interface{}) (interface{}, error) {
	req := request.(*proto.CreateUserRequest)
	dest := &commands.CreateUserCommand{}
	err := mapstructure.Decode(req, dest)
	if err != nil {
		return nil, err
	}
	return dest, nil
}
