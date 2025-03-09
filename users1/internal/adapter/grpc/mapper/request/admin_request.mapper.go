package request

import (
	"context"

	"github.com/baothaihcmut/Ecommerce-Go/users/internal/adapter/grpc/proto"
	valueobject "github.com/baothaihcmut/Ecommerce-Go/users/internal/core/command/domain/aggregates/user/value_object"
	"github.com/baothaihcmut/Ecommerce-Go/users/internal/core/command/port/inbound/commands"
)

type AdminRequestMapper interface {
	ToAdminLoginCommand(_ context.Context, request interface{}) (interface{}, error)
	ToAdminVerifyTokenCommand(_ context.Context, request interface{}) (interface{}, error)
}

type AdminRequestMapperImpl struct{}

func (a *AdminRequestMapperImpl) ToAdminLoginCommand(_ context.Context, request interface{}) (interface{}, error) {
	req := request.(*proto.AdminLoginRequest)
	email, err := valueobject.NewEmail(req.Email)
	if err != nil {
		return nil, err
	}
	return &commands.LoginCommand{
		Email:    *email,
		Password: req.Password,
	}, nil
}

func (a *AdminRequestMapperImpl) ToAdminVerifyTokenCommand(_ context.Context, request interface{}) (interface{}, error) {
	req := request.(*proto.AdminVerifyTokenRequest)
	return &commands.VerifyTokenCommand{
		Token:          req.Token,
		IsRefreshToken: req.IsRefreshToken,
	}, nil
}

func NewAdminRequestMapper() AdminRequestMapper {
	return &AdminRequestMapperImpl{}
}
