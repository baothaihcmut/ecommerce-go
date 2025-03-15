package commands

import (
	valueobjects "github.com/baothaihcmut/Ecommerce-go/products/internal/core/command/domain/aggregates/products/value_objects"
)

type UpdateProductCommand struct {
	Id          valueobjects.ProductId
	Name        *string
	Description *string
	Unit        *string
}
