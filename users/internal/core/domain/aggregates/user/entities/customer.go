package entities

import (
	valueobject "github.com/baothaihcmut/Ecommerce-Go/users/internal/core/domain/aggregates/user/value_object"
)

type Customer struct {
	LoyaltyPoint valueobject.LoyaltyPoint
	Rank         valueobject.Rank
}

func NewCustomer() (*Customer, error) {
	point, err := valueobject.NewLoyaltyPoin(0)
	if err != nil {
		return nil, err
	}
	return &Customer{
		LoyaltyPoint: *point,
		Rank:         valueobject.NewRank(*point),
	}, nil
}

func NewCustomerWithPoint(point valueobject.LoyaltyPoint) *Customer {
	return &Customer{
		LoyaltyPoint: point,
		Rank:         valueobject.NewRank(point),
	}
}
