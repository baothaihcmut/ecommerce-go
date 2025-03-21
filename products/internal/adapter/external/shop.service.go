package external

import (
	"context"

	"github.com/baothaihcmut/Ecommerce-go/products/internal/core/port/outbound/external/models"
	"github.com/google/uuid"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type ShopService struct {
	inMemoryShops []*models.Shop
}

func NewShopService() *ShopService {
	shopID, _ := primitive.ObjectIDFromHex("65f01a4e9b3e2c7d4f1a8b9c")
	ownerId, _ := uuid.FromBytes([]byte("c6047b65-d8b0-4899-ae27-88c6c1ecdbd6"))
	return &ShopService{
		inMemoryShops: []*models.Shop{
			{
				ID:          shopID,
				ShopOwnerId: ownerId,
			},
		},
	}
}
func (s *ShopService) FindShopById(ctx context.Context, id primitive.ObjectID) (*models.Shop, error) {
	for _, shop := range s.inMemoryShops {
		if shop.ID == id {
			return shop, nil
		}
	}
	return nil, nil
}
