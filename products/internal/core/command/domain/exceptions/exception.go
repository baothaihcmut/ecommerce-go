package exceptions

import "errors"

var (
	ErrDuplicateVariation          = errors.New("duplicate variation in one product")
	ErrMismatchVariationValue      = errors.New("you cannot add this varition value to this item")
	ErrVariationExist              = errors.New("variation exist in product")
	ErrCategoryExist               = errors.New("category exist in product")
	ErrVariationNotBelongToProduct = errors.New("variation not belong to product")
	ErrPriceLessThanZero           = errors.New("product price cannot be less than 0")
	ErrProductQuantityLessThanZero = errors.New("product quantity cannot be less than 0")
)
