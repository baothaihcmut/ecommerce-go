package services

import (
	"context"
	"fmt"

	v1 "github.com/baothaihcmut/Ecommerce-go/libs/pkg/proto/shared/v1"
	userProto "github.com/baothaihcmut/Ecommerce-go/libs/pkg/proto/users/v1"
	"github.com/baothaihcmut/Ecommerce-go/users/internal/adapter/grpc/exception"
	"github.com/baothaihcmut/Ecommerce-go/users/internal/adapter/grpc/mappers"
	"github.com/baothaihcmut/Ecommerce-go/users/internal/core/port/inbound/handlers"
	"google.golang.org/grpc/status"
)

type AuthService struct {
	userProto.UnsafeAuthServiceServer
	authHandler handlers.AuthHandler
}

// ConfirmSignUp implements v1.AuthServiceServer.
func (a *AuthService) ConfirmSignUp(context.Context, *userProto.ConfirmSignUpRequest) (*userProto.ConfirmSignUpResponse, error) {
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
func (a *AuthService) SignUp(ctx context.Context, req *userProto.SignUpRequest) (*userProto.SignUpResponse, error) {
	res, err := a.authHandler.SignUp(ctx, mappers.ToSignUpCommand(req))
	if err != nil {
		fmt.Println(err)
		msg, code := exception.MapException(err)
		return &userProto.SignUpResponse{
			Data: nil,
			Status: &v1.Status{
				Message: msg,
				Success: false,
			},
		}, status.Error(code, msg)
	}
	return &userProto.SignUpResponse{
		Data: mappers.ToSignUpResponse(res),
		Status: &v1.Status{
			Success: true,
			Message: "Sign up success",
		},
	}, nil
}

func NewAuthService(authHandler handlers.AuthHandler) userProto.AuthServiceServer {
	return &AuthService{
		authHandler: authHandler,
	}
}
