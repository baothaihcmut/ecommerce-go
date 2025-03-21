package services

import (
	"context"
	"net/http"

	v1 "github.com/baothaihcmut/Ecommerce-go/libs/pkg/proto/shared/v1"
	userProto "github.com/baothaihcmut/Ecommerce-go/libs/pkg/proto/users/v1"
	"github.com/baothaihcmut/Ecommerce-go/users/internal/adapter/grpc/mappers"
	"github.com/baothaihcmut/Ecommerce-go/users/internal/core/port/inbound/handlers"
)

type AuthService struct {
	userProto.UnimplementedAuthServiceServer
	authHandler handlers.AuthHandler
}

// ConfirmSignUp implements v1.AuthServiceServer.
func (a *AuthService) ConfirmSignUp(ctx context.Context, req *userProto.ConfirmSignUpRequest) (*userProto.ConfirmSignUpResponse, error) {

	res, err := a.authHandler.ConfirmSignUp(ctx, mappers.ToConfirmSignUpCommand(req))
	if err != nil {
		return nil, err
	}
	return &userProto.ConfirmSignUpResponse{
		Data: mappers.ToConfirmSignUpResponse(res),
		Status: &v1.Status{
			Message: "Confirm sign up success",
			Code:    http.StatusCreated,
		},
	}, nil
}

// LogIn implements v1.AuthServiceServer.
func (a *AuthService) LogIn(ctx context.Context, req *userProto.LogInRequest) (*userProto.LogInResponse, error) {
	res, err := a.authHandler.LogIn(ctx, mappers.ToLogInCommand(req))
	if err != nil {
		return nil, err
	}

	return &userProto.LogInResponse{
		Data: &userProto.LogInData{
			AccessToken:  res.AccessToken,
			RefreshToken: res.RefreshToken,
		},
		Status: &v1.Status{
			Message: "Log in success",
			Code:    http.StatusCreated,
		},
	}, nil
}
func (a *AuthService) SignUp(ctx context.Context, req *userProto.SignUpRequest) (*userProto.SignUpResponse, error) {
	res, err := a.authHandler.SignUp(ctx, mappers.ToSignUpCommand(req))
	if err != nil {
		return nil, err
	}
	return &userProto.SignUpResponse{
		Data: mappers.ToSignUpResponse(res),
		Status: &v1.Status{
			Code:    http.StatusCreated,
			Message: "Sign up success",
		},
	}, nil
}

func NewAuthService(authHandler handlers.AuthHandler) userProto.AuthServiceServer {
	return &AuthService{
		authHandler: authHandler,
	}
}
