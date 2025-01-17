package repositories

import (
	"context"

	"github.com/baothaihcmut/Ecommerce-Go/products/internal/core/domain/aggregates/shops"
	valueobjects "github.com/baothaihcmut/Ecommerce-Go/products/internal/core/domain/aggregates/shops/value_objects"
)

type ShopRepository interface {
	FindShopById(context.Context, valueobjects.ShopId) (shops.Shop, error)
}
