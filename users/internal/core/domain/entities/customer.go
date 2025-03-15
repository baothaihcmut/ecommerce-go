package entities

import "github.com/baothaihcmut/Ecommerce-go/users/internal/core/domain/enums"

type Customer struct {
	User         *User
	LoyaltyPoint int
	Rank         enums.CustomerRank
}

func NewCustomer(
	user *User,
) *Customer {
	return &Customer{
		User:         user,
		LoyaltyPoint: 0,
		Rank:         enums.CUSTOMER_RANK_BRONZE,
	}
}
