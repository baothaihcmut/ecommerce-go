package services

import (
	"database/sql"
	"errors"

	"github.com/baothaihcmut/Ecommerce-Go/users/internal/core/port/inbound/command/handlers"
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

func NewUserCommandService(userRepo outbound.UserRepository, dbSource *sql.DB) handlers.UserCommandHandler {
	return &UserCommandService{
		userRepo: userRepo,
		dbSource: dbSource,
	}
}
