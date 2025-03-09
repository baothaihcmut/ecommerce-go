package response

import (
	"context"

	"github.com/baothaihcmut/Ecommerce-Go/users/internal/adapter/grpc/proto"
	"github.com/baothaihcmut/Ecommerce-Go/users/internal/core/command/port/inbound/results"
	"github.com/google/uuid"
)

type AdminResponseMapper interface {
	ToAdminLoginReponse(_ context.Context, resp interface{}) (interface{}, error)
	ToAdminVerifyTokenResponse(_ context.Context, resp interface{}) (interface{}, error)
}

type AdminResponseMapperImpl struct {
}

func (a *AdminResponseMapperImpl) ToAdminLoginReponse(_ context.Context, resp interface{}) (interface{}, error) {
	res := resp.(*results.LoginCommandResult)
	return &proto.AdminLoginData{
		AccessToken:  res.AccessToken,
		RefreshToken: res.RefreshToken,
	}, nil
}

func (a *AdminResponseMapperImpl) ToAdminVerifyTokenResponse(_ context.Context, resp interface{}) (interface{}, error) {
	res := resp.(*results.VerifyTokenCommandResult)
	return &proto.AdminVerifyTokenData{
		Id: uuid.UUID(res.Id).String(),
	}, nil
}
func NewAdminResponseMapper() AdminResponseMapper {
	return &AdminResponseMapperImpl{}
}
