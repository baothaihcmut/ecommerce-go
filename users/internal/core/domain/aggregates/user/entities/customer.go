package entities

import (
	valueobject "github.com/baothaihcmut/Ecommerce-Go/users/internal/core/domain/aggregates/user/value_object"
)

type Customer struct {
	Id           valueobject.UserId
	LoyaltyPoint valueobject.LoyaltyPoint
	Rank         valueobject.Rank
}

func NewCustomer(Id valueobject.UserId) (*Customer, error) {
	point, err := valueobject.NewLoyaltyPoin(0)
	if err != nil {
		return nil, err
	}
	return &Customer{
		Id:           Id,
		LoyaltyPoint: *point,
		Rank:         valueobject.NewRank(*point),
	}, nil
}

func NewCustomerWithPoint(Id valueobject.UserId, point valueobject.LoyaltyPoint) *Customer {
	return &Customer{
		Id:           Id,
		LoyaltyPoint: point,
		Rank:         valueobject.NewRank(point),
	}
}
