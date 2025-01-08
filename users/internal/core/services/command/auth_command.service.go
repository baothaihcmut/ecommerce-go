package services

import (
	"context"
	"errors"

	"github.com/baothaihcmut/Ecommerce-Go/users/internal/core/domain/aggregates/user"
	"github.com/baothaihcmut/Ecommerce-Go/users/internal/core/port/inbound/command/commands"
	"github.com/baothaihcmut/Ecommerce-Go/users/internal/core/port/inbound/command/handlers"
	"github.com/baothaihcmut/Ecommerce-Go/users/internal/core/port/inbound/command/results"
	"github.com/baothaihcmut/Ecommerce-Go/users/internal/core/port/outbound"
)

var (
	ErrInvalidToken = errors.New("invalid token")
	ErrTokenExpire  = errors.New("token expire")
)

type AuthCommandService struct {
	jwtPort  outbound.JwtPort
	userRepo outbound.UserRepository
}

func (s *AuthCommandService) Login(ctx context.Context, command *commands.LoginCommand) (*results.LoginCommandResult, error) {
	//get user from db
	userDb, err := s.userRepo.FindByEmail(ctx, command.Email)
	//
	if err != nil {
		return nil, err
	}
	//check if user exist
	if userDb == nil {
		return nil, user.ErrBadCredencial
	}
	//validate password
	err = userDb.ValidateAuth(command.Password)
	if err != nil {
		return nil, err
	}
	//generate access token and refresh token
	accessToken, err := s.jwtPort.GenerateAccessToken(ctx, userDb)
	if err != nil {
		return nil, err
	}
	refreshToken, err := s.jwtPort.GenerateRefreshToken(ctx, userDb)
	if err != nil {
		return nil, err
	}
	//return access and refresh token
	return &results.LoginCommandResult{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	}, nil
}

func NewAuthCommandService(userRepo outbound.UserRepository, jwtPort outbound.JwtPort) handlers.AuthCommandHandler {
	return &AuthCommandService{
		userRepo: userRepo,
		jwtPort:  jwtPort,
	}
}
