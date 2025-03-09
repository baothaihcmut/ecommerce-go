package services

import (
	"context"
	"errors"

	"github.com/baothaihcmut/Ecommerce-Go/libs/pkg/postgres"
	"github.com/baothaihcmut/Ecommerce-Go/libs/pkg/tracing"
	"github.com/baothaihcmut/Ecommerce-Go/users/internal/core/command/domain/aggregates/user"
	valueobject "github.com/baothaihcmut/Ecommerce-Go/users/internal/core/command/domain/aggregates/user/value_object"
	"github.com/baothaihcmut/Ecommerce-Go/users/internal/core/command/domain/enums"
	"github.com/baothaihcmut/Ecommerce-Go/users/internal/core/command/port/inbound/commands"
	"github.com/baothaihcmut/Ecommerce-Go/users/internal/core/command/port/inbound/handlers"
	"github.com/baothaihcmut/Ecommerce-Go/users/internal/core/command/port/inbound/results"
	"github.com/baothaihcmut/Ecommerce-Go/users/internal/core/command/port/outbound"
	"github.com/google/uuid"
	"go.opentelemetry.io/otel/trace"
)

var (
	ErrInvalidToken = errors.New("invalid token")
	ErrTokenExpire  = errors.New("token expire")
)

type AuthCommandService struct {
	jwtPort   outbound.JwtService
	userRepo  outbound.UserRepository
	dbService postgres.TransactionService
	tracer    trace.Tracer
}

// VerifyToken implements handlers.AuthCommandHandler.

func (s *AuthCommandService) Login(ctx context.Context, command *commands.LoginCommand) (_ *results.LoginCommandResult, err error) {
	ctx, span := tracing.StartSpan(ctx, s.tracer, "Auth.Login: service", map[string]interface{}{
		"email": string(command.Email),
	})
	defer tracing.EndSpan(span, err, nil)
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
	accessToken, err := s.jwtPort.GenerateAccessToken(ctx, outbound.GenerateTokenArg{
		UserId: uuid.UUID(userDb.Id),
		Role:   userDb.Role,
	})
	if err != nil {
		return nil, err
	}
	refreshToken, err := s.jwtPort.GenerateRefreshToken(ctx, outbound.GenerateTokenArg{
		UserId: uuid.UUID(userDb.Id),
		Role:   userDb.Role,
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

func (s *AuthCommandService) toUserDomain(command *commands.SignUpCommand) (*user.User, error) {
	email, err := valueobject.NewEmail(command.Email)
	if err != nil {
		return nil, err
	}
	phoneNumber, err := valueobject.NewPhoneNumber(command.PhoneNumber)
	if err != nil {
		return nil, err
	}
	password, err := valueobject.NewPassword(command.Password)
	if err != nil {
		return nil, err
	}

	if command.Role == enums.CUSTOMER {
		return user.NewCustomer(
			*email, password, *phoneNumber, command.Addresses, command.FirstName, command.LastName,
		)
	} else {
		return user.NewShopOwner(
			*email, password, *phoneNumber, command.Addresses, command.FirstName, command.LastName, command.ShopOwnerInfo.BussinessLincese,
		)
	}
}

func (s *AuthCommandService) SignUp(ctx context.Context, command *commands.SignUpCommand) (_ *results.SignUpCommandResult, err error) {
	ctx, span := tracing.StartSpan(ctx, s.tracer, "Auth.SignUp: service", map[string]interface{}{
		"email": string(command.Email),
	})
	defer tracing.EndSpan(span, err, nil)
	user, err := s.toUserDomain(command)
	if err != nil {
		return nil, err
	}
	//check if email exist
	emailExist, err := s.userRepo.CheckEmailExist(ctx, user.Email)
	if err != nil {
		return nil, err
	}
	if emailExist {
		return nil, ErrEmailExist
	}

	//check if phone number exist
	phoneExist, err := s.userRepo.CheckPhoneNumberExist(ctx, user.PhoneNumber)
	if err != nil {
		return nil, err
	}
	if phoneExist {
		return nil, ErrPhoneNumberExist
	}
	//generate token
	accessToken, err := s.jwtPort.GenerateAccessToken(ctx, outbound.GenerateTokenArg{
		UserId: uuid.UUID(user.Id),
		Role:   user.Role,
	})
	if err != nil {
		return nil, err
	}
	refreshToken, err := s.jwtPort.GenerateRefreshToken(ctx, outbound.GenerateTokenArg{
		UserId: uuid.UUID(user.Id),
		Role:   user.Role,
	})
	if err != nil {
		return nil, err
	}

	err = user.SetCurrentRefreshToken(valueobject.NewToken(refreshToken, enums.REFRESH_TOKEN))
	if err != nil {
		return nil, err
	}
	//persist to db
	tx, err := s.dbService.BeginTransaction(ctx)
	defer func() { _ = s.dbService.RollbackTransaction(ctx, tx) }()
	if err != nil {
		return nil, err
	}
	err = s.userRepo.Save(ctx, user, tx)
	if err != nil {
		return nil, err
	}
	return &results.SignUpCommandResult{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	}, s.dbService.CommitTransaction(ctx, tx)

}
func (s *AuthCommandService) VerifyToken(ctx context.Context, command *commands.VerifyTokenCommand) (_ *results.VerifyTokenCommandResult, err error) {
	ctx, span := tracing.StartSpan(ctx, s.tracer, "Auth.VerifyToken: service", nil)
	defer tracing.EndSpan(span, err, nil)
	if !command.IsRefreshToken {
		accesToken, err := s.jwtPort.DecodeAccessToken(ctx, command.Token)
		if err != nil {
			return nil, err
		}
		return &results.VerifyTokenCommandResult{
			Id:   accesToken.Id,
			Role: accesToken.Role,
		}, nil
	} else {
		refreshToken, err := s.jwtPort.DecodeRefreshToken(ctx, command.Token)
		if err != nil {
			return nil, err
		}
		return &results.VerifyTokenCommandResult{
			Id: refreshToken.Id,
		}, nil
	}
}

func NewAuthCommandService(userRepo outbound.UserRepository, jwtPort outbound.JwtService, dbService postgres.TransactionService, tracer trace.Tracer) handlers.AuthCommandHandler {
	return &AuthCommandService{
		userRepo:  userRepo,
		jwtPort:   jwtPort,
		dbService: dbService,
		tracer:    tracer,
	}
}
