package request

import (
	"context"

	"github.com/baothaihcmut/Ecommerce-Go/users/internal/adapter/grpc/proto"
	valueobject "github.com/baothaihcmut/Ecommerce-Go/users/internal/core/domain/aggregates/user/value_object"
	"github.com/baothaihcmut/Ecommerce-Go/users/internal/core/port/inbound/command/commands"
)

type AuthRequestMapper interface {
	ToLoginCommand(context.Context, interface{}) (interface{}, error)
}

type AuthRequestMapperImpl struct {
}

func (m *AuthRequestMapperImpl) ToLoginCommand(_ context.Context, request interface{}) (interface{}, error) {
	req := request.(*proto.LoginRequest)
	email, err := valueobject.NewEmail(req.Email)
	if err != nil {
		return nil, err
	}
	return &commands.LoginCommand{
		Email:    *email,
		Password: req.Password,
	}, nil
}
