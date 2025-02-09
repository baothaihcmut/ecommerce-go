package services

import (
	"context"

	"github.com/baothaihcmut/Ecommerce-Go/users/internal/core/domain/aggregates/user"
	"github.com/baothaihcmut/Ecommerce-Go/users/internal/core/domain/enums"
	"github.com/baothaihcmut/Ecommerce-Go/users/internal/core/port/inbound/command/commands"
	"github.com/baothaihcmut/Ecommerce-Go/users/internal/core/port/inbound/command/handlers"
	"github.com/baothaihcmut/Ecommerce-Go/users/internal/core/port/inbound/command/results"
	"github.com/baothaihcmut/Ecommerce-Go/users/internal/core/port/outbound"
	"github.com/google/uuid"
)

type AdminCommandService struct {
	adminRepo outbound.AdminRepository
	jwtPort   outbound.JwtService
}

func (a *AdminCommandService) LogIn(ctx context.Context, command *commands.LoginCommand) (*results.LoginCommandResult, error) {
	admin, err := a.adminRepo.FindByEmail(ctx, command.Email)
	if err != nil {
		return nil, err
	}
	if admin == nil {
		return nil, user.ErrBadCredencial
	}
	err = admin.LogIn(command.Password)
	if err != nil {
		return nil, err
	}
	//generate access token and refresh token
	accessToken, err := a.jwtPort.GenerateAccessToken(ctx, outbound.GenerateTokenArg{
		UserId: uuid.UUID(admin.Id),
		Role:   enums.ADMIN,
	})
	if err != nil {
		return nil, err
	}
	refreshToken, err := a.jwtPort.GenerateRefreshToken(ctx, outbound.GenerateTokenArg{
		UserId: uuid.UUID(admin.Id),
		Role:   enums.ADMIN,
	})
	if err != nil {
		return nil, err
	}
	//return access and refresh token
	return &results.LoginCommandResult{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	}, nil
}
func (a *AdminCommandService) VerifyToken(ctx context.Context, command *commands.VerifyTokenCommand) (*results.VerifyTokenCommandResult, error) {
	if !command.IsRefreshToken {
		accesToken, err := a.jwtPort.DecodeAccessToken(ctx, command.Token)
		if err != nil {
			return nil, err
		}
		return &results.VerifyTokenCommandResult{
			Id:   accesToken.Id,
			Role: accesToken.Role,
		}, nil
	} else {
		refreshToken, err := a.jwtPort.DecodeRefreshToken(ctx, command.Token)
		if err != nil {
			return nil, err
		}
		return &results.VerifyTokenCommandResult{
			Id: refreshToken.Id,
		}, nil
	}
}

func NewAdminCommandService(adminRepo outbound.AdminRepository, jwtService outbound.JwtService) handlers.AdminCommandHandler {
	return &AdminCommandService{
		adminRepo: adminRepo,
		jwtPort:   jwtService,
	}
}
