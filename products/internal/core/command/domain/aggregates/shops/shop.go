package shops

import valueobjects "github.com/baothaihcmut/Ecommerce-Go/products/internal/core/command/domain/aggregates/shops/value_objects"

type Shop struct {
	Id       valueobjects.ShopId
	IsActive bool
}
