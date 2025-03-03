package handlers

import (
	"context"

	"github.com/baothaihcmut/Ecommerce-Go/users/internal/core/command/port/inbound/commands"
	"github.com/baothaihcmut/Ecommerce-Go/users/internal/core/command/port/inbound/results"
)

type AuthCommandHandler interface {
	Login(context.Context, *commands.LoginCommand) (*results.LoginCommandResult, error)
	SignUp(context.Context, *commands.SignUpCommand) (*results.SignUpCommandResult, error)
	VerifyToken(context.Context, *commands.VerifyTokenCommand) (*results.VerifyTokenCommandResult, error)
}
