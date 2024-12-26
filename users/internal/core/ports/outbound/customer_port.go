package outbound

import (
	"github.com/baothaihcmut/Ecommerce-Go/users/internal/core/domain/aggregates/customer"
	valueobject "github.com/baothaihcmut/Ecommerce-Go/users/internal/core/domain/aggregates/customer/value_object"
)

type CustomerPort interface {
	Save(*customer.Customer) error
	FindById(valueobject.UserId) (*customer.Customer, error)
}
