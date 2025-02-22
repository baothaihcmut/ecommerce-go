package services

import (
	"context"
	"errors"

	"github.com/baothaihcmut/Ecommerce-Go/libs/pkg/postgres"
	"github.com/baothaihcmut/Ecommerce-Go/libs/pkg/tracing"
	"github.com/baothaihcmut/Ecommerce-Go/users/internal/core/command/domain/aggregates/user"
	valueobject "github.com/baothaihcmut/Ecommerce-Go/users/internal/core/command/domain/aggregates/user/value_object"
	"github.com/baothaihcmut/Ecommerce-Go/users/internal/core/command/domain/enums"
	"github.com/baothaihcmut/Ecommerce-Go/users/internal/core/command/events"
	"github.com/baothaihcmut/Ecommerce-Go/users/internal/core/command/exception"
	"github.com/baothaihcmut/Ecommerce-Go/users/internal/core/command/port/inbound/commands"
	"github.com/baothaihcmut/Ecommerce-Go/users/internal/core/command/port/inbound/handlers"
	"github.com/baothaihcmut/Ecommerce-Go/users/internal/core/command/port/inbound/results"
	"github.com/baothaihcmut/Ecommerce-Go/users/internal/core/command/port/outbound/external"
	"github.com/baothaihcmut/Ecommerce-Go/users/internal/core/command/port/outbound/repositories"
	"github.com/google/uuid"
	"go.opentelemetry.io/otel/trace"
)

var (
	ErrInvalidToken = errors.New("invalid token")
	ErrTokenExpire  = errors.New("token expire")
)

type AuthCommandService struct {
	jwtPort            external.JwtService
	userRepo           repositories.UserRepository
	userConfirmService external.UserConfirmService
	dbService          postgres.TransactionService
	tracer             trace.Tracer
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
	accessToken, err := s.jwtPort.GenerateAccessToken(ctx, external.GenerateAccessTokenArg{
		UserId:            uuid.UUID(userDb.Id),
		IsShopOwnerActive: userDb.IsShopOwnerActive,
	})
	if err != nil {
		return nil, err
	}
	refreshToken, err := s.jwtPort.GenerateRefreshToken(ctx, external.GenerateRefreshTokenArg{
		UserId: uuid.UUID(userDb.Id),
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
	addrArgs := make([]user.AddressArg, len(command.Addresses))
	for idx, addr := range command.Addresses {
		addrArgs[idx] = user.AddressArg{
			Street: addr.Street,
			City:   addr.City,
			Town:   addr.Town,
		}
	}
	return user.NewUser(
		*email,
		password,
		*phoneNumber,
		addrArgs,
		command.FirstName,
		command.LastName,
	)
}

func (s *AuthCommandService) SignUp(ctx context.Context, command *commands.SignUpCommand) (_ *results.SignUpCommandResult, err error) {
	ctx, span := tracing.StartSpan(ctx, s.tracer, "Auth.SignUp: service", map[string]interface{}{
		"email": string(command.Email),
	})
	defer tracing.EndSpan(span, err, nil)
	newUser, err := s.toUserDomain(command)
	if err != nil {
		return nil, err
	}
	//check if email exist
	emailExist, err := s.userRepo.CheckEmailExist(ctx, newUser.Email)
	if err != nil {
		return nil, err
	}
	if emailExist {
		return nil, exception.ErrEmailExist
	}

	//check if phone number exist
	phoneExist, err := s.userRepo.CheckPhoneNumberExist(ctx, newUser.PhoneNumber)
	if err != nil {
		return nil, err
	}
	if phoneExist {
		return nil, exception.ErrPhoneNumberExist
	}
	//if user activate shopowner
	if command.IsShopOwnerActive {
		newUser.ActivateShopOwner(user.ActivateShopOwnerArg{
			BussinessLincese: command.ShopOwnerInfo.BussinessLincese,
		})
	}
	//store user in pending
	code, err := s.userConfirmService.StoreUserInfo(ctx, newUser)
	if err != nil {
		return nil, err
	}
	//publish event for mail sending,...
	e := events.UserSignUpEvent{
		Email:          string(newUser.Email),
		PhoneNumber:    string(newUser.PhoneNumber),
		FirstName:      newUser.FirstName,
		LastName:       newUser.LastName,
		VericationCode: code,
	}
	if err = s.userConfirmService.PublishSignUpEvent(ctx, &e); err != nil {
		return nil, err
	}

	return &results.SignUpCommandResult{}, nil

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
			Role: enums.USER,
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

func NewAuthCommandService(userRepo repositories.UserRepository, jwtPort external.JwtService, dbService postgres.TransactionService, tracer trace.Tracer) handlers.AuthCommandHandler {
	return &AuthCommandService{
		userRepo:  userRepo,
		jwtPort:   jwtPort,
		dbService: dbService,
		tracer:    tracer,
	}
}
