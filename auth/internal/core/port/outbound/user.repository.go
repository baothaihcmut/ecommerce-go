package outbound

import (
	"context"

	"github.com/baothaihcmut/Ecommerce-Go/auth/internal/core/domain"
	valueobject "github.com/baothaihcmut/Ecommerce-Go/auth/internal/core/domain/value_object"
)

type UserRepository interface {
	FindUserByEmail(context.Context, valueobject.Email) (*domain.User, error)
	Save(context.Context, *domain.User) error
}
