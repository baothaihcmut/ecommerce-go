package repositories

import (
	"context"

	"github.com/baothaihcmut/Ecommerce-Go/products/internal/core/command/domain/aggregates/shops"
	valueobjects "github.com/baothaihcmut/Ecommerce-Go/products/internal/core/command/domain/aggregates/shops/value_objects"
)

type ShopCommandRepository interface {
	FindShopById(context.Context, valueobjects.ShopId) (shops.Shop, error)
}
