package repositories

import (
	"context"

	"github.com/baothaihcmut/Ecommerce-Go/products/internal/core/domain/aggregates/products"
	"go.mongodb.org/mongo-driver/mongo"
)

type ProductRepository interface {
	Save(context.Context, *products.Product, mongo.Session) error
}
