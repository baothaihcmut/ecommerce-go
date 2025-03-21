package repositories

import (
	"context"

	"github.com/baothaihcmut/Ecommerce-go/products/internal/core/domain/entities"
)

type ProductRepo interface {
	CreateProduct(context.Context, *entities.Product) error
}
