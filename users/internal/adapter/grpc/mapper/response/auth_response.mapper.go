package response

import (
	"context"

	"github.com/baothaihcmut/Ecommerce-Go/users/internal/adapter/grpc/proto"
	commandResult "github.com/baothaihcmut/Ecommerce-Go/users/internal/core/port/inbound/command/results"
)

type AuthResponseMapper interface {
	ToLoginResult(context.Context, interface{}) (interface{}, error)
}

type AuthResponseMapperImpl struct{}

func (m *AuthResponseMapperImpl) ToLoginResult(_ context.Context, res interface{}) (interface{}, error) {
	resp := res.(commandResult.LoginCommandResult)
	return &proto.LoginData{
		AccessToken:  resp.AccessToken.Value,
		RefreshToken: resp.RefreshToken.Value,
	}, nil
}