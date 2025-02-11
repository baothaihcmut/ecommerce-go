package handlers

import (
	"context"

	"github.com/baothaihcmut/Ecommerce-Go/users/internal/core/command/port/inbound/commands"
	"github.com/baothaihcmut/Ecommerce-Go/users/internal/core/command/port/inbound/results"
)

type AdminCommandHandler interface {
	LogIn(context.Context, *commands.LoginCommand) (*results.LoginCommandResult, error)
	VerifyToken(context.Context, *commands.VerifyTokenCommand) (*results.VerifyTokenCommandResult, error)
}
