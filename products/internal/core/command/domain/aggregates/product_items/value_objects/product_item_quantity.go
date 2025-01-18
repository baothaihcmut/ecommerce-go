package valueobjects

import "errors"

var (
	ErrProductQuantityLessThanZero = errors.New("product quantity cannot be less than 0")
)

type ProductQuantity int

func NewProductQuantity(quantity int) (ProductQuantity, error) {
	if quantity < 0 {
		return ProductQuantity(0), ErrProductQuantityLessThanZero
	}
	return ProductQuantity(quantity), nil
}
func (p ProductQuantity) IncreQuantity(price ProductQuantity) ProductQuantity {
	return ProductQuantity(int(p) + int(price))
}

func (p ProductQuantity) DecreQuantity(price ProductQuantity) (ProductQuantity, error) {
	if int(p) < int(price) {
		return ProductQuantity(0), ErrProductQuantityLessThanZero
	}
	return ProductQuantity(int(p) - int(price)), nil
}
