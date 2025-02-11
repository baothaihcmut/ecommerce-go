package commands

import valueobject "github.com/baothaihcmut/Ecommerce-Go/users/internal/core/command/domain/aggregates/user/value_object"

type LoginCommand struct {
	Email    valueobject.Email
	Password string
}
