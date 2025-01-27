package commands

import (
	categoryValueobjects "github.com/baothaihcmut/Ecommerce-Go/products/internal/core/command/domain/aggregates/categories/value_objects"
	productValueobjects "github.com/baothaihcmut/Ecommerce-Go/products/internal/core/command/domain/aggregates/products/value_objects"
)

type AddProductCategoiesCommand struct {
	ProductId     productValueobjects.ProductId
	NewCategories []categoryValueobjects.CategoryId
}
