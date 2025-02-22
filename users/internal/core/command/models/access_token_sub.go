package models

import (
	valueobject "github.com/baothaihcmut/Ecommerce-Go/users/internal/core/command/domain/aggregates/user/value_object"
)

type AccessTokenSub struct {
	Id                valueobject.UserId
	IsShopOwnerActive bool
}
