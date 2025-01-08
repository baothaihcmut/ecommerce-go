package results

import (
	"github.com/baothaihcmut/Ecommerce-Go/users/internal/core/domain/aggregates/user/entities"
	valueobject "github.com/baothaihcmut/Ecommerce-Go/users/internal/core/domain/aggregates/user/value_object"
	"github.com/baothaihcmut/Ecommerce-Go/users/internal/core/domain/enums"
)

type CreateUserCommandResult struct {
	Id          valueobject.UserId
	Email       valueobject.Email
	PhoneNumber valueobject.PhoneNumber
	Address     []valueobject.Address
	Role        enums.Role
	FirstName   string
	LastName    string
	Customer    *entities.Customer
	ShopOwner   *entities.ShopOwner
}
