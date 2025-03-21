package external

import (
	"context"

	"github.com/baothaihcmut/Ecommerce-go/products/internal/core/port/outbound/external/models"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type ShopService interface {
	FindShopById(context.Context, primitive.ObjectID) (*models.Shop, error)
}
