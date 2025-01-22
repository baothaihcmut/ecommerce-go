package results

import (
	"github.com/baothaihcmut/Ecommerce-Go/products/internal/core/command/domain/aggregates/products"
)

type CreateProductResult struct {
	*products.Product
}
