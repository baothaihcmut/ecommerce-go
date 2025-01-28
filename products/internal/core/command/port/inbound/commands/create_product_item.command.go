package commands

import (
	productitems "github.com/baothaihcmut/Ecommerce-Go/products/internal/core/command/domain/aggregates/product_items"
	productItemValueobjects "github.com/baothaihcmut/Ecommerce-Go/products/internal/core/command/domain/aggregates/product_items/value_objects"
	productValueobjects "github.com/baothaihcmut/Ecommerce-Go/products/internal/core/command/domain/aggregates/products/value_objects"
)

type ItemVariation struct {
	Variation string
	Value     string
}

type CreateProductItemCommand struct {
	Sku             productItemValueobjects.SKU
	ProductId       productValueobjects.ProductId
	Price           productItemValueobjects.ProductPrice
	Quantity        productItemValueobjects.ProductQuantity
	VariationValues []productitems.VariationValueArg
}
