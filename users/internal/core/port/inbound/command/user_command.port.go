package inbound

import (
	"context"

	"github.com/baothaihcmut/Ecommerce-Go/users/internal/core/commands"
	result "github.com/baothaihcmut/Ecommerce-Go/users/internal/core/results/command"
)

type UserCommandPort interface {
	CreateUser(context.Context, *commands.CreateUserCommand) (*result.CreateUserCommandResult, error)
}
