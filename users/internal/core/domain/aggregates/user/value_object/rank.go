package valueobject

import "github.com/baothaihcmut/Ecommerce-Go/users/internal/core/domain/enums"

type Rank enums.CustomerRank

func NewRank(p LoyaltyPoint) Rank {
	if p < 50 {
		return Rank(enums.BRONZE)
	} else if p >= 50 && p <= 100 {
		return Rank((enums.SILVER))
	} else {
		return Rank(enums.GOLD)
	}
}
