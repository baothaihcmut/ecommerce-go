package outbound

import (
	"context"

	"github.com/baothaihcmut/Ecommerce-Go/users/internal/core/domain/aggregates/user"
	valueobject "github.com/baothaihcmut/Ecommerce-Go/users/internal/core/domain/aggregates/user/value_object"
	"github.com/baothaihcmut/Ecommerce-Go/users/internal/core/domain/models"
)

type JwtPort interface {
	DecodeAccessToken(context.Context, valueobject.Token) (models.AccessTokenSub, error)
	DecodeRefreshToken(context.Context, valueobject.Token) (models.RefreshTokenSub, error)
	GenerateAccessToken(context.Context, *user.User) (valueobject.Token, error)
	GenerateRefreshToken(context.Context, *user.User) (valueobject.Token, error)
}
