package result

import (
	valueobject "github.com/baothaihcmut/Ecommerce-Go/users/internal/core/domain/aggregates/user/value_object"
	"github.com/baothaihcmut/Ecommerce-Go/users/internal/core/domain/enums"
)

type ValidateUserResult struct {
	Id   valueobject.UserId
	Role enums.Role
}
