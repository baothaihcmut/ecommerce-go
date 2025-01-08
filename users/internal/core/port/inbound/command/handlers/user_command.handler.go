package handlers

import (
	"context"

	"github.com/baothaihcmut/Ecommerce-Go/users/internal/core/port/inbound/command/commands"
	"github.com/baothaihcmut/Ecommerce-Go/users/internal/core/port/inbound/command/results"
)

type UserCommandHandler interface {
	CreateUser(context.Context, *commands.CreateUserCommand) (*results.CreateUserCommandResult, error)
}
