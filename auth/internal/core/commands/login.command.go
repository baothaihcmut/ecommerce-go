package commands

import valueobject "github.com/baothaihcmut/Ecommerce-Go/auth/internal/core/domain/value_object"

type LoginCommand struct {
	GoogleToken valueobject.Token
}
