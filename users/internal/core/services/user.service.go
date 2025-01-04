package services

import (
	"context"
	"database/sql"
	"errors"

	"github.com/baothaihcmut/Ecommerce-Go/users/internal/core/commands"
	"github.com/baothaihcmut/Ecommerce-Go/users/internal/core/domain/aggregates/user"
	valueobject "github.com/baothaihcmut/Ecommerce-Go/users/internal/core/domain/aggregates/user/value_object"
	"github.com/baothaihcmut/Ecommerce-Go/users/internal/core/domain/enums"
	"github.com/baothaihcmut/Ecommerce-Go/users/internal/core/port/inbound"
	"github.com/baothaihcmut/Ecommerce-Go/users/internal/core/port/outbound"
	"github.com/baothaihcmut/Ecommerce-Go/users/internal/core/queries"
)

var (
	ErrEmailExist       = errors.New("email exist")
	ErrPhoneNumberExist = errors.New("phone number exist")
)

type UserService struct {
	userRepo outbound.UserRepository
	dbSource *sql.DB
}

func (u *UserService) toUserDomain(command *commands.CreateUserCommand) (*user.User, error) {
	email, err := valueobject.NewEmail(command.Email)
	if err != nil {
		return nil, err
	}
	phoneNumber, err := valueobject.NewPhoneNumber(command.PhoneNumber)
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
			*email, *phoneNumber, address, command.FirstName, command.LastName,
		)
	} else {
		return user.NewShopOwner(
			*email, *phoneNumber, address, command.FirstName, command.LastName, command.ShopOwnerInfo.BussinessLincese,
		)
	}
}

// CreateUser implements inbound.UserService.
func (u *UserService) CreateUser(ctx context.Context, command *commands.CreateUserCommand) (*user.User, error) {
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
		return user, tx.Commit()
	case err := <-errCh:
		return nil, err
	case <-ctx.Done():
		return nil, ctx.Err()
	}
}

// FindUserById implements inbound.UserService.
func (u *UserService) FindUserById(ctx context.Context, q *queries.FindUserByIdQuery) (*user.User, error) {
	userId, err := valueobject.NewUserId(q.Id)
	if err != nil {
		return nil, err
	}
	doneCh := make(chan *user.User)
	errCh := make(chan error)
	go func() {
		user, err := u.userRepo.FindById(ctx, *userId)
		if err != nil {
			errCh <- err
		}
		doneCh <- user
	}()
	select {
	case user := <-doneCh:
		return user, nil
	case err := <-errCh:
		return nil, err
	case <-ctx.Done():
		return nil, ctx.Err()
	}
}

func NewUserService(userRepo outbound.UserRepository, dbSource *sql.DB) inbound.UserService {
	return &UserService{
		userRepo: userRepo,
		dbSource: dbSource,
	}
}
