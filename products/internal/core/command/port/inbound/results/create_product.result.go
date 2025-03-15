package results

import (
	"github.com/baothaihcmut/Ecommerce-go/products/internal/core/command/domain/aggregates/products"
)

type CreateProductResult struct {
	*products.Product
}
