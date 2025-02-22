package valueobject

import "github.com/baothaihcmut/Ecommerce-Go/users/internal/core/command/exception"

type LoyaltyPoint int16

func NewLoyaltyPoin(point int16) (*LoyaltyPoint, error) {
	if point < 0 {
		return nil, exception.InvalidPoint
	}
	return (*LoyaltyPoint)(&point), nil
}
