package outbound

import (
	"context"

	"github.com/baothaihcmut/Ecommerce-Go/auth/internal/core/domain"
	valueobject "github.com/baothaihcmut/Ecommerce-Go/auth/internal/core/domain/value_object"
)

type JwtRepository interface {
	GenerateAccessToken(context.Context, *domain.User) (valueobject.Token, error)
	GenerateRefreshToken(context.Context, *domain.User) (valueobject.Token, error)
	DecodeToken(context.Context, valueobject.Token) (*domain.User, error)
}
