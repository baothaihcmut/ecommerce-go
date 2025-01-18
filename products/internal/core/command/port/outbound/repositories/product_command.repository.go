package repositories

import (
	"context"

	"github.com/baothaihcmut/Ecommerce-Go/products/internal/core/command/domain/aggregates/products"
	"go.mongodb.org/mongo-driver/mongo"
)

type ProductCommandRepository interface {
	Save(context.Context, *products.Product, mongo.Session) error
}
