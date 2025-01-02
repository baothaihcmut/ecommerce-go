package utils

import (
	valueobject "github.com/baothaihcmut/Ecommerce-Go/users/internal/core/domain/aggregates/user/value_object"
	"github.com/baothaihcmut/Ecommerce-Go/users/internal/core/services"
	"google.golang.org/grpc/codes"
)

func MapErrorToGrpcStatus(err error) codes.Code {
	switch {
	case err == valueobject.InvalidEmail ||
		err == valueobject.InvalidPhonenumber ||
		err == valueobject.InvalidPoint:
		return codes.InvalidArgument
	case err == services.ErrEmailExist ||
		err == services.ErrPhoneNumberExist:
		return codes.AlreadyExists
	default:
		return codes.Internal
	}
}
