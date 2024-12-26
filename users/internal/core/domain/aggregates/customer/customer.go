package customer

import (
	"github.com/baothaihcmut/Ecommerce-Go/users/internal/core/domain/aggregates/customer/entities"
	valueobject "github.com/baothaihcmut/Ecommerce-Go/users/internal/core/domain/aggregates/customer/value_object"
	"github.com/google/uuid"
)

type Customer struct {
	Id           valueobject.UserId
	User         *entities.User
	LoyaltyPoint valueobject.LoyaltyPoint
	Rank         valueobject.Rank
}

func NewCustomer(
	user *entities.User,
) (*Customer, error) {
	point, err := valueobject.NewLoyaltyPoin(0)
	if err != nil {
		return nil, err
	}
	return &Customer{
		Id:           valueobject.UserId(uuid.New()),
		User:         user,
		LoyaltyPoint: *point,
		Rank:         valueobject.NewRank(*point),
	}, nil
}
