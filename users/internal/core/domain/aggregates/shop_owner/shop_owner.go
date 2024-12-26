package shopowner

import (
	"github.com/baothaihcmut/Ecommerce-Go/users/internal/core/domain/aggregates/customer/entities"
	valueobject "github.com/baothaihcmut/Ecommerce-Go/users/internal/core/domain/aggregates/customer/value_object"
)

type ShopOwner struct {
	Id               valueobject.UserId
	User             *entities.User
	BussinessLicense string
}
