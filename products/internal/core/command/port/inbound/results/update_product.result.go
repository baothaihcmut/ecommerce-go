package results

import "github.com/baothaihcmut/Ecommerce-Go/products/internal/core/command/domain/aggregates/products"

type UpdateProductResult struct {
	*products.Product
}
