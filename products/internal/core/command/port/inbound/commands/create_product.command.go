package commands

import (
	categoryValueobjects "github.com/baothaihcmut/Ecommerce-go/products/internal/core/command/domain/aggregates/categories/value_objects"
	valueobjects "github.com/baothaihcmut/Ecommerce-go/products/internal/core/command/domain/aggregates/products/value_objects"
)

type CreateProductCommand struct {
	Name        string
	Description string
	Unit        string
	ShopId      valueobjects.ShopId
	CategoryIds []categoryValueobjects.CategoryId
	Variations  []string
	Images      []CreateImageCommand
}
