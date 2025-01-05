package inbound

import (
	"context"

	"github.com/baothaihcmut/Ecommerce-Go/users/internal/core/queries"
	result "github.com/baothaihcmut/Ecommerce-Go/users/internal/core/results/query"
)

type UserQueryPort interface {
	FindUserById(context.Context, *queries.FindUserByIdQuery) (*result.UserQueryResult, error)
}
