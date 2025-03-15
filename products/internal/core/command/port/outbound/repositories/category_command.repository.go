package repositories

import (
	"context"

	"github.com/baothaihcmut/Ecommerce-go/products/internal/core/command/domain/aggregates/categories"
	valueobjects "github.com/baothaihcmut/Ecommerce-go/products/internal/core/command/domain/aggregates/categories/value_objects"
	"go.mongodb.org/mongo-driver/mongo"
)

type CategoryCommandRepository interface {
	Save(ctx context.Context, category *categories.Category, session mongo.Session) error
	BulkSave(ctx context.Context, categories []*categories.Category, session mongo.Session) error
	FindCategoryById(ctx context.Context, categoryId valueobjects.CategoryId) (*categories.Category, error)
}
