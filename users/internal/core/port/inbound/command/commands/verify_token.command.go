package commands

import valueobject "github.com/baothaihcmut/Ecommerce-Go/users/internal/core/domain/aggregates/user/value_object"

type VerifyTokenCommand struct {
	Token valueobject.Token
}
