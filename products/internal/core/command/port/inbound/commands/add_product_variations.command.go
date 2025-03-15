package commands

import (
	productValueobjects "github.com/baothaihcmut/Ecommerce-go/products/internal/core/command/domain/aggregates/products/value_objects"
)

type AddProductVariationsCommand struct {
	ProductId     productValueobjects.ProductId
	NewVariations []string
}
