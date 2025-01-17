package repositories

import (
	"context"

	"github.com/baothaihcmut/Ecommerce-Go/products/internal/core/domain/aggregates/categories"
	valueobjects "github.com/baothaihcmut/Ecommerce-Go/products/internal/core/domain/aggregates/categories/value_objects"
	"go.mongodb.org/mongo-driver/mongo"
)

type CategoryRepository interface {
	Save(context.Context, *categories.Category, mongo.Session) error
	FindCategoryById(ctx context.Context, categoryId valueobjects.CategoryId) (*categories.Category, error)
}
