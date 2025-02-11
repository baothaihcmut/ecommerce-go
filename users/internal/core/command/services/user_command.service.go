package services

import (
	"database/sql"
	"errors"

	"github.com/baothaihcmut/Ecommerce-Go/users/internal/core/command/port/inbound/handlers"
	"github.com/baothaihcmut/Ecommerce-Go/users/internal/core/command/port/outbound"
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

func NewUserCommandService(userRepo outbound.UserRepository, dbSource *sql.DB) handlers.UserCommandHandler {
	return &UserCommandService{
		userRepo: userRepo,
		dbSource: dbSource,
	}
}
