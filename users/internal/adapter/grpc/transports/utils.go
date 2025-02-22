package transports

import (
	"github.com/baothaihcmut/Ecommerce-Go/users/internal/core/command/domain/aggregates/user"
	valueobject "github.com/baothaihcmut/Ecommerce-Go/users/internal/core/command/domain/aggregates/user/value_object"
	"github.com/baothaihcmut/Ecommerce-Go/users/internal/core/command/exception"
	"google.golang.org/grpc/codes"
)

func MapErrorToGrpcStatus(err error) codes.Code {
	switch {
	case err == valueobject.InvalidEmail ||
		err == exception.InvalidPhonenumber ||
		err == exception.InvalidPoint:
		return codes.InvalidArgument
	case err == exception.ErrEmailExist ||
		err == exception.ErrPhoneNumberExist:
		return codes.AlreadyExists
	case err == user.ErrMisMatchRefreshToken:
		return codes.Unauthenticated
	default:
		return codes.Internal
	}
}
