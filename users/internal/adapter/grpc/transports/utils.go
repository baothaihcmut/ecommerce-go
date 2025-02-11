package transports

import (
	"github.com/baothaihcmut/Ecommerce-Go/users/internal/core/command/domain/aggregates/user"
	valueobject "github.com/baothaihcmut/Ecommerce-Go/users/internal/core/command/domain/aggregates/user/value_object"
	commandServices "github.com/baothaihcmut/Ecommerce-Go/users/internal/core/command/services"
	"google.golang.org/grpc/codes"
)

func MapErrorToGrpcStatus(err error) codes.Code {
	switch {
	case err == valueobject.InvalidEmail ||
		err == valueobject.InvalidPhonenumber ||
		err == valueobject.InvalidPoint:
		return codes.InvalidArgument
	case err == commandServices.ErrEmailExist ||
		err == commandServices.ErrPhoneNumberExist:
		return codes.AlreadyExists
	case err == user.ErrMisMatchRefreshToken:
		return codes.Unauthenticated
	default:
		return codes.Internal
	}
}
