package handlers

import (
	"context"

	"github.com/baothaihcmut/Ecommerce-go/users/internal/core/port/inbound/commands"
	"github.com/baothaihcmut/Ecommerce-go/users/internal/core/port/inbound/results"
)

type AuthHandler interface {
	SignUp(context.Context, *commands.SignUpCommand) (*results.SignUpResult, error)
}
