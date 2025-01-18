package valueobjects

import (
	"github.com/baothaihcmut/Ecommerce-Go/products/internal/core/command/domain/exceptions"
)

type ProductPrice int

func NewProductPrice(price int) (ProductPrice, error) {
	if price < 0 {
		return ProductPrice(0), exceptions.ErrPriceLessThanZero
	}
	return ProductPrice(price), nil
}

func (p ProductPrice) IncrePrice(price ProductPrice) ProductPrice {
	return ProductPrice(int(p) + int(price))
}

func (p ProductPrice) DecrePrice(price ProductPrice) (ProductPrice, error) {
	if int(p) < int(price) {
		return ProductPrice(0), exceptions.ErrPriceLessThanZero
	}
	return ProductPrice(int(p) - int(price)), nil
}
