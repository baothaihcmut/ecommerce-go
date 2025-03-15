package valueobjects

import (
	"github.com/baothaihcmut/Ecommerce-go/products/internal/core/command/domain/exceptions"
)

type ProductQuantity int

func NewProductQuantity(quantity int) (ProductQuantity, error) {
	if quantity < 0 {
		return ProductQuantity(0), exceptions.ErrProductQuantityLessThanZero
	}
	return ProductQuantity(quantity), nil
}
func (p ProductQuantity) IncreQuantity(price ProductQuantity) ProductQuantity {
	return ProductQuantity(int(p) + int(price))
}

func (p ProductQuantity) DecreQuantity(price ProductQuantity) (ProductQuantity, error) {
	if int(p) < int(price) {
		return ProductQuantity(0), exceptions.ErrProductQuantityLessThanZero
	}
	return ProductQuantity(int(p) - int(price)), nil
}
