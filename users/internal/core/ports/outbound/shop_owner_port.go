package outbound

import (
	valueobject "github.com/baothaihcmut/Ecommerce-Go/users/internal/core/domain/aggregates/customer/value_object"
	shopowner "github.com/baothaihcmut/Ecommerce-Go/users/internal/core/domain/aggregates/shop_owner"
)

type ShopOwnerPort interface {
	Save(*shopowner.ShopOwner) error
	FindById(valueobject.UserId) error
}
