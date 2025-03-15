package errors

import (
	domainException "github.com/baothaihcmut/Ecommerce-go/products/internal/core/command/domain/exceptions"
	commandException "github.com/baothaihcmut/Ecommerce-go/products/internal/core/command/exceptions"
	"google.golang.org/grpc/codes"
)

func MapGrpcErrorCode(err error) codes.Code {
	switch err {
	case domainException.ErrCategoryExist,
		domainException.ErrDuplicateVariation:
		return codes.AlreadyExists
	case domainException.ErrMismatchVariationValue,
		domainException.ErrPriceLessThanZero,
		domainException.ErrVariationNotBelongToProduct,
		domainException.ErrMismatchVariationValue,
		domainException.ErrProductQuantityLessThanZero,
		commandException.ErrVariationOfItemNotBelongToProduct:
		return codes.InvalidArgument
	case commandException.ErrParentCategoryNotExist, commandException.ErrProductNotExist:
		return codes.NotFound
	default:
		return codes.Internal
	}
}
