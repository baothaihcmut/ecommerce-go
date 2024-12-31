package inbound

import (
	"context"

	"github.com/baothaihcmut/Ecommerce-Go/users/internal/core/commands"
	"github.com/baothaihcmut/Ecommerce-Go/users/internal/core/domain/aggregates/user"
	"github.com/baothaihcmut/Ecommerce-Go/users/internal/core/queries"
)

type UserService interface {
	CreateUser(context.Context, *commands.CreateUserCommand) (*user.User, error)
	FindUserById(context.Context, *queries.FindUserByIdQuery) (*user.User, error)
}
