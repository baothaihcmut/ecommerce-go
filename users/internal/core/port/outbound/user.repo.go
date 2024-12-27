package outbound

import (
	"context"
	"database/sql"

	"github.com/baothaihcmut/Ecommerce-Go/users/internal/core/domain/aggregates/user"
	valueobject "github.com/baothaihcmut/Ecommerce-Go/users/internal/core/domain/aggregates/user/value_object"
)

type UserRepository interface {
	Save(context.Context, *user.User, *sql.Tx) error
	FindById(context.Context, valueobject.UserId) (*user.User, error)
}
