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
	ErrEmailExist             = errors.New("email exist")
	ErrPhoneNumberExist       = errors.New("phone number exist")
	ErrInvalidEmailOrPassword = errors.New("user name or password incorrect")
)

type UserCommandService struct {
	userRepo outbound.UserRepository
	dbSource *sql.DB
}

func (u *UserCommandService) toUserDomain(command *commands.CreateUserCommand) (*user.User, error) {
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

// CreateUser implements inbound.UserCommandService.
func (u *UserCommandService) CreateUser(ctx context.Context, command *commands.CreateUserCommand) (*results.CreateUserCommandResult, error) {
	user, err := u.toUserDomain(command)
	if err != nil {
		return nil, err
	}
	//check if email exist
	emailExist, err := u.userRepo.FindByEmail(ctx, user.Email)
	if err != nil && err != sql.ErrNoRows {
		return nil, err
	}
	if emailExist != nil {
		return nil, ErrEmailExist
	}

	//check if phone number exist
	phoneExist, err := u.userRepo.FindByPhoneNumber(ctx, user.PhoneNumber)
	if err != nil && err != sql.ErrNoRows {
		return nil, err
	}
	if phoneExist != nil {
		return nil, ErrPhoneNumberExist
	}
	doneCh := make(chan interface{})
	errCh := make(chan error)
	tx, err := u.dbSource.Begin()
	defer tx.Rollback()
	if err != nil {
		return nil, err
	}
	go func() {
		err = u.userRepo.Save(ctx, user, tx)
		if err != nil {
			errCh <- err
		}
		doneCh <- 1
	}()

	select {
	case <-doneCh:
		return &results.CreateUserCommandResult{
			Id:          user.Id,
			Email:       user.Email,
			PhoneNumber: user.PhoneNumber,
			FirstName:   user.FirstName,
			LastName:    user.LastName,
			Customer:    user.Customer,
			ShopOwner:   user.ShopOwner,
		}, tx.Commit()
	case err := <-errCh:
		return nil, err
	case <-ctx.Done():
		return nil, ctx.Err()
	}
}

// FindUserById implements inbound.UserCommandService.

func NewUserCommandService(userRepo outbound.UserRepository, dbSource *sql.DB) handlers.UserCommandHandler {
	return &UserCommandService{
		userRepo: userRepo,
		dbSource: dbSource,
	}
}

// ValidateUser implements inbound.UserQueryPort.
