package inmemory

import (
	"context"

	"github.com/baothaihcmut/Ecommerce-Go/products/internal/core/command/domain/aggregates/shops"
	valueobjects "github.com/baothaihcmut/Ecommerce-Go/products/internal/core/command/domain/aggregates/shops/value_objects"
	"github.com/baothaihcmut/Ecommerce-Go/products/internal/core/command/port/outbound/repositories"
)

type InMemoryShopRepository struct {
	storage map[string]*shops.Shop
}

// FindShopById implements repositories.ShopRepository.
func (i *InMemoryShopRepository) FindShopById(ctx context.Context, shopId valueobjects.ShopId) (shops.Shop, error) {
	return *i.storage[string(shopId)], nil
}

func NewInMemoryShopRepository() repositories.ShopCommandRepository {
	return &InMemoryShopRepository{
		storage: map[string]*shops.Shop{
			"64e5f2b6c529fb27c83647e2": {
				Id:       "64e5f2b6c529fb27c83647e2",
				IsActive: true,
			},
			"64e5f2b6c529fb27c83647e3": {
				Id:       "64e5f2b6c529fb27c83647e3",
				IsActive: false,
			},
		},
	}
}
