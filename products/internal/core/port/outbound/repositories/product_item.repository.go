package repositories

import (
	"context"

	productitems "github.com/baothaihcmut/Ecommerce-Go/products/internal/core/domain/aggregates/product_items"
	"go.mongodb.org/mongo-driver/mongo"
)

type ProductItemRepository interface {
	Save(context.Context, *productitems.ProductItem, mongo.Session) error
}
