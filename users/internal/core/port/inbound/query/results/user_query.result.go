package results

import (
	"github.com/baothaihcmut/Ecommerce-Go/users/internal/core/domain/aggregates/user/entities"
	valueobject "github.com/baothaihcmut/Ecommerce-Go/users/internal/core/domain/aggregates/user/value_object"
	"github.com/baothaihcmut/Ecommerce-Go/users/internal/core/domain/enums"
)

type UserQueryResult struct {
	Id          valueobject.UserId
	Email       valueobject.Email
	PhoneNumber valueobject.PhoneNumber
	Address     []*entities.Address
	Role        enums.Role
	FirstName   string
	LastName    string
	Customer    *entities.Customer
	ShopOwner   *entities.ShopOwner
}
