package handlers

import (
	"context"

	"github.com/baothaihcmut/Ecommerce-Go/users/internal/core/port/inbound/query/queries"
	"github.com/baothaihcmut/Ecommerce-Go/users/internal/core/port/inbound/query/results"
)

type UserQueryHandler interface {
	FindUserById(context.Context, *queries.FindUserByIdQuery) (*results.UserQueryResult, error)
}
