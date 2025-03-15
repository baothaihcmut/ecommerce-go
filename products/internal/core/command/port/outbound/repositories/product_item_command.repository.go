package repositories

import (
	"context"

	productitems "github.com/baothaihcmut/Ecommerce-go/products/internal/core/command/domain/aggregates/product_items"
	"go.mongodb.org/mongo-driver/mongo"
)

type ProductItemCommandRepository interface {
	Save(context.Context, *productitems.ProductItem, mongo.Session) error
}
