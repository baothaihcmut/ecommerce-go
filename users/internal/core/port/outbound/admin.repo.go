package outbound

import (
	"context"

	"github.com/baothaihcmut/Ecommerce-Go/users/internal/core/domain/aggregates/admin"
	valueobject "github.com/baothaihcmut/Ecommerce-Go/users/internal/core/domain/aggregates/user/value_object"
)

type AdminRepository interface {
	Save(context.Context, *admin.Admin) error
	FindByEmail(context.Context, valueobject.Email) (*admin.Admin, error)
}
