package inbound

import (
	"context"

	"github.com/baothaihcmut/Ecommerce-Go/auth/internal/core/commands"
	"github.com/baothaihcmut/Ecommerce-Go/auth/internal/core/results"
)

type AuthInboundPort interface {
	Login(context.Context, *commands.LoginCommand) (*results.LoginResult, error)
}
