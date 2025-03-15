package inmemory

import (
	"context"

	"github.com/baothaihcmut/Ecommerce-go/products/internal/core/command/models"
	"github.com/baothaihcmut/Ecommerce-go/products/internal/core/command/services"
)

type InMemoryShopRepository struct {
	storage map[string]*models.Shop
}

// FindShopById implements repositories.ShopRepository.
func (i *InMemoryShopRepository) FindShopById(ctx context.Context, shopId string) (*models.Shop, error) {
	return i.storage[string(shopId)], nil
}

func NewInMemoryShopService() services.ShopService {
	return &InMemoryShopRepository{
		storage: map[string]*models.Shop{
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
