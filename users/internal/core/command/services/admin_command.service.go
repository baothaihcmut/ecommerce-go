package services

import (
	"context"

	"github.com/baothaihcmut/Ecommerce-Go/libs/pkg/tracing"
	"github.com/baothaihcmut/Ecommerce-Go/users/internal/core/command/domain/aggregates/user"
	"github.com/baothaihcmut/Ecommerce-Go/users/internal/core/command/domain/enums"
	"github.com/baothaihcmut/Ecommerce-Go/users/internal/core/command/port/inbound/commands"
	"github.com/baothaihcmut/Ecommerce-Go/users/internal/core/command/port/inbound/handlers"
	"github.com/baothaihcmut/Ecommerce-Go/users/internal/core/command/port/inbound/results"
	"github.com/baothaihcmut/Ecommerce-Go/users/internal/core/command/port/outbound"
	"github.com/google/uuid"
	"go.opentelemetry.io/otel/trace"
)

type AdminCommandService struct {
	adminRepo outbound.AdminRepository
	jwtPort   outbound.JwtService
	tracer    trace.Tracer
}

func (a *AdminCommandService) LogIn(ctx context.Context, command *commands.LoginCommand) (res *results.LoginCommandResult, err error) {
	ctx, span := tracing.StartSpan(ctx, a.tracer, "Admin.LogIn: service", map[string]interface{}{
		"email": string(command.Email),
	})
	defer tracing.EndSpan(span, err, nil)
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
func (a *AdminCommandService) VerifyToken(ctx context.Context, command *commands.VerifyTokenCommand) (res *results.VerifyTokenCommandResult, err error) {
	ctx, span := tracing.StartSpan(ctx, a.tracer, "Admin.LogIn: service", nil)
	defer tracing.EndSpan(span, err, nil)
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

func NewAdminCommandService(adminRepo outbound.AdminRepository, jwtService outbound.JwtService, tracer trace.Tracer) handlers.AdminCommandHandler {
	return &AdminCommandService{
		adminRepo: adminRepo,
		jwtPort:   jwtService,
		tracer:    tracer,
	}
}
