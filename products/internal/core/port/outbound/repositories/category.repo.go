package repositories

import (
	"context"

	"github.com/baothaihcmut/Ecommerce-go/products/internal/core/domain/entities"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type CategoryRepo interface {
	FindCategoryById(context.Context, primitive.ObjectID) (*entities.Category, error)
}
