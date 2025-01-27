package repositories

import (
	"context"

	"github.com/baothaihcmut/Ecommerce-Go/products/internal/core/command/domain/aggregates/products"
	valueobjects "github.com/baothaihcmut/Ecommerce-Go/products/internal/core/command/domain/aggregates/products/value_objects"
	"go.mongodb.org/mongo-driver/mongo"
)

type ProductCommandRepository interface {
	Save(context.Context, *products.Product, mongo.Session) error
	FindById(context.Context, valueobjects.ProductId) (*products.Product, error)
}
