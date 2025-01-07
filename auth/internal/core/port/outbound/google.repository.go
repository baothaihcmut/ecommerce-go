package outbound

import (
	"context"

	valueobject "github.com/baothaihcmut/Ecommerce-Go/auth/internal/core/domain/value_object"
)

type OAuthRepository interface {
	GetUserEmail(context.Context, valueobject.Token) (valueobject.Email, error)
}
