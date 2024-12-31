package response

import (
	"context"

	"github.com/baothaihcmut/Ecommerce-Go/users/internal/adapter/grpc/proto"
	"github.com/baothaihcmut/Ecommerce-Go/users/internal/core/domain/aggregates/user"
	"github.com/mitchellh/mapstructure"
)

type UserResponseMapper interface {
	ToCreateUserResponse(_ context.Context, response interface{}) (interface{}, error)
}

type UserResponseMapperImpl struct {
}

func (m *UserResponseMapperImpl) ToCreateUserResponse(_ context.Context, response interface{}) (interface{}, error) {
	res := response.(*user.User)
	dest := &proto.CreateUserResponse{}
	err := mapstructure.Decode(res, dest)
	if err != nil {
		return nil, err
	}
	return dest, nil
}
