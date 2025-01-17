package repositories

import (
	"context"

	"github.com/baothaihcmut/Ecommerce-Go/products/internal/core/domain/aggregates/categories"
	"go.mongodb.org/mongo-driver/mongo"
)

type CategoryRepository interface {
	Save(context.Context, *categories.Category, mongo.Session) error
}
