package results

import valueobject "github.com/baothaihcmut/Ecommerce-Go/users/internal/core/domain/aggregates/user/value_object"

type LoginCommandResult struct {
	AccessToken  valueobject.Token
	RefreshToken valueobject.Token
}
