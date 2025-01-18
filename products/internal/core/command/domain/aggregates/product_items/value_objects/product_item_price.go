package valueobjects

import "errors"

type ProductPrice int

var (
	ErrPriceLessThanZero = errors.New("product price cannot be less than 0")
)

func NewProductPrice(price int) (ProductPrice, error) {
	if price < 0 {
		return ProductPrice(0), ErrPriceLessThanZero
	}
	return ProductPrice(price), nil
}

func (p ProductPrice) IncrePrice(price ProductPrice) ProductPrice {
	return ProductPrice(int(p) + int(price))
}

func (p ProductPrice) DecrePrice(price ProductPrice) (ProductPrice, error) {
	if int(p) < int(price) {
		return ProductPrice(0), ErrPriceLessThanZero
	}
	return ProductPrice(int(p) - int(price)), nil
}
