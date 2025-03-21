package results

import (
	"time"

	"github.com/baothaihcmut/Ecommerce-go/products/internal/core/domain/entities"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type ProductResult struct {
	ID          primitive.ObjectID
	Name        string
	Description string
	CategoryIds []primitive.ObjectID
	ShopId      primitive.ObjectID
	Variations  []string
	SoldTotal   int
	CreatedAt   time.Time
	UpdatedAt   time.Time
}

func MapToProductResult(product *entities.Product) ProductResult {
	return ProductResult{
		ID:          product.ID,
		Name:        product.Name,
		Description: product.Description,
		CategoryIds: product.CategoryIds,
		ShopId:      product.ShopId,
		Variations:  product.Variations,
		SoldTotal:   product.SoldTotal,
		CreatedAt:   product.CreatedAt,
		UpdatedAt:   product.UpdatedAt,
	}
}
