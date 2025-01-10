package results

import (
	valueobject "github.com/baothaihcmut/Ecommerce-Go/users/internal/core/domain/aggregates/user/value_object"
)

type SignUpCommandResult struct {
	AccessToken  valueobject.Token
	RefreshToken valueobject.Token
}
