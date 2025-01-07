package results

import valueobject "github.com/baothaihcmut/Ecommerce-Go/auth/internal/core/domain/value_object"

type LoginResult struct {
	AccessToken  valueobject.Token
	RefreshToken valueobject.Token
}
