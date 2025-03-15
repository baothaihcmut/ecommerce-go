package services

import (
	"context"

	userProto "github.com/baothaihcmut/Ecommerce-go/libs/pkg/proto/users/v1"
)

type AuthService struct {
	userProto.UnsafeAuthServiceServer
}

// LogIn implements v1.AuthServiceServer.
func (a *AuthService) LogIn(context.Context, *userProto.LogInRequest) (*userProto.LogInResponse, error) {
	return &userProto.LogInResponse{
		Data: &userProto.LogInData{
			AccessToken:  "hello",
			RefreshToken: "hello",
		},
	}, nil
}

func NewAuthService() userProto.AuthServiceServer {
	return &AuthService{}
}
