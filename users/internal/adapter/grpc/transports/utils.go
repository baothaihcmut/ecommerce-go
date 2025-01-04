package transports

import (
	valueobject "github.com/baothaihcmut/Ecommerce-Go/users/internal/core/domain/aggregates/user/value_object"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func mapErrorToGrpcStatus(err error) error {
	switch {
	case err == valueobject.InvalidEmail || err == valueobject.InvalidPhonenumber || err == valueobject.InvalidPoint:
		return status.Error(codes.InvalidArgument, err.Error())
	}
}
