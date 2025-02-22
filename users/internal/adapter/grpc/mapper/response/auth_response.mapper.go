package response

import (
	"context"

	"github.com/baothaihcmut/Ecommerce-Go/users/internal/adapter/grpc/proto"
	commandResult "github.com/baothaihcmut/Ecommerce-Go/users/internal/core/command/port/inbound/results"
	"github.com/google/uuid"
)

type AuthResponseMapper interface {
	ToLoginResult(context.Context, interface{}) (interface{}, error)
	ToSignUpResult(_ context.Context, res interface{}) (interface{}, error)
	ToVerifyTokenResult(_ context.Context, res interface{}) (interface{}, error)
}

type AuthResponseMapperImpl struct{}

func (m *AuthResponseMapperImpl) ToLoginResult(_ context.Context, res interface{}) (interface{}, error) {
	resp := res.(*commandResult.LoginCommandResult)
	return &proto.LoginData{
		AccessToken:  resp.AccessToken,
		RefreshToken: resp.RefreshToken,
	}, nil
}

func (m *AuthResponseMapperImpl) ToSignUpResult(_ context.Context, res interface{}) (interface{}, error) {
	resp := res.(*commandResult.SignUpCommandResult)
	return &proto.LoginData{
		AccessToken:  resp.AccessToken,
		RefreshToken: resp.RefreshToken,
	}, nil
}

func (m *AuthResponseMapperImpl) ToVerifyTokenResult(_ context.Context, res interface{}) (interface{}, error) {
	resp := res.(*commandResult.VerifyTokenCommandResult)
	return &proto.VerifyTokenData{
		Id: uuid.UUID(resp.Id).String(),
	}, nil
}
func NewAuthResponseMapper() AuthResponseMapper {
	return &AuthResponseMapperImpl{}
}
