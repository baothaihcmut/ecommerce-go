package valueobject

import "errors"

type LoyaltyPoint int16

var (
	InvalidPoint = errors.New("Point must greater than 0")
)

func NewLoyaltyPoin(point int16) (*LoyaltyPoint, error) {
	if point < 0 {
		return nil, InvalidPoint
	}
	return (*LoyaltyPoint)(&point), nil
}
