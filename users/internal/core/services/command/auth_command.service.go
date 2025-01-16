package services

import (
	"context"
	"database/sql"
	"errors"

	"github.com/baothaihcmut/Ecommerce-Go/users/internal/core/domain/aggregates/user"
	valueobject "github.com/baothaihcmut/Ecommerce-Go/users/internal/core/domain/aggregates/user/value_object"
	"github.com/baothaihcmut/Ecommerce-Go/users/internal/core/domain/enums"
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
	dbSource *sql.DB
}

// VerifyToken implements handlers.AuthCommandHandler.

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
	address := make([]valueobject.Address, len(command.Addresses))
	for idx, addr := range command.Addresses {
		address[idx] = *valueobject.NewAddress(
			addr.Priority, addr.Street, addr.Town, addr.City, addr.Province,
		)
	}

	if command.Role == enums.CUSTOMER {
		return user.NewCustomer(
			*email, password, *phoneNumber, address, command.FirstName, command.LastName,
		)
	} else {
		return user.NewShopOwner(
			*email, password, *phoneNumber, address, command.FirstName, command.LastName, command.ShopOwnerInfo.BussinessLincese,
		)
	}
}

func (s *AuthCommandService) SignUp(ctx context.Context, command *commands.SignUpCommand) (*results.SignUpCommandResult, error) {
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
	accessToken, err := s.jwtPort.GenerateAccessToken(ctx, user)
	if err != nil {
		return nil, err
	}
	refreshToken, err := s.jwtPort.GenerateRefreshToken(ctx, user)
	if err != nil {
		return nil, err
	}
	user.SetCurrentRefreshToken(refreshToken)
	//persist to db
	tx, err := s.dbSource.Begin()
	defer tx.Rollback()
	if err != nil {
		return nil, err
	}
	s.userRepo.Save(ctx, user, tx)
	return &results.SignUpCommandResult{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	}, tx.Commit()

}
func (s *AuthCommandService) VerifyToken(ctx context.Context, command *commands.VerifyTokenCommand) (*results.VerifyTokenCommandResult, error) {
	if command.Token.TokenType == enums.ACCESS_TOKEN {
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

func NewAuthCommandService(userRepo outbound.UserRepository, jwtPort outbound.JwtPort, dbSource *sql.DB) handlers.AuthCommandHandler {
	return &AuthCommandService{
		userRepo: userRepo,
		jwtPort:  jwtPort,
		dbSource: dbSource,
	}
}
