package models

import (
	valueobject "github.com/baothaihcmut/Ecommerce-Go/users/internal/core/command/domain/aggregates/user/value_object"
	"github.com/baothaihcmut/Ecommerce-Go/users/internal/core/command/domain/enums"
)

type AccessTokenSub struct {
	Id   valueobject.UserId
	Role enums.Role
}
