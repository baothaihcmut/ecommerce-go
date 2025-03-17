package entities

import "github.com/baothaihcmut/Ecommerce-go/users/internal/core/domain/enums"

type Customer struct {
	LoyaltyPoint int
	Rank         enums.CustomerRank
}

func NewCustomer() *Customer {
	return &Customer{
		LoyaltyPoint: 0,
		Rank:         enums.CUSTOMER_RANK_BRONZE,
	}
}
