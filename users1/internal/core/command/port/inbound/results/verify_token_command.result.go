package results

import (
	valueobject "github.com/baothaihcmut/Ecommerce-Go/users/internal/core/command/domain/aggregates/user/value_object"
	"github.com/baothaihcmut/Ecommerce-Go/users/internal/core/command/domain/enums"
)

type VerifyTokenCommandResult struct {
	Id   valueobject.UserId
	Role enums.Role
}
