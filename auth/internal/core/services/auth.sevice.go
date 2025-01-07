package services

import (
	"context"
	"errors"

	"github.com/baothaihcmut/Ecommerce-Go/auth/internal/core/commands"
	"github.com/baothaihcmut/Ecommerce-Go/auth/internal/core/port/outbound"
	"github.com/baothaihcmut/Ecommerce-Go/auth/internal/core/results"
)

var (
	ErrUserNotRegister = errors.New("User email doesn't exist")
)

type AuthService struct {
	oauthRepository outbound.OAuthRepository
	jwtRepository   outbound.JwtRepository
	userRepository  outbound.UserRepository
}

func (s *AuthService) Login(ctx context.Context, command *commands.LoginCommand) (*results.LoginResult, error) {
	//get user email by oath repo
	email, err := s.oauthRepository.GetUserEmail(ctx, command.GoogleToken)
	if err != nil {
		return nil, err
	}
	//get user from db by email
	user, err := s.userRepository.FindUserByEmail(ctx, email)
	if err != nil {
		return nil, err
	}
	// if user not found throw error
	if user == nil {
		return nil, ErrUserNotRegister
	}
	//generate token
	accessToken, err := s.jwtRepository.GenerateAccessToken(ctx, user)
	if err != nil {
		return nil, err
	}
	refreshToken, err := s.jwtRepository.GenerateRefreshToken(ctx, user)
	if err != nil {
		return nil, err
	}
	return &results.LoginResult{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	}, nil
}
