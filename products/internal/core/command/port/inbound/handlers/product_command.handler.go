package handlers

import (
	"context"

	"github.com/baothaihcmut/Ecommerce-go/products/internal/core/command/port/inbound/commands"
	"github.com/baothaihcmut/Ecommerce-go/products/internal/core/command/port/inbound/results"
)

type ProductCommandHandler interface {
	CreateProduct(ctx context.Context, product *commands.CreateProductCommand) (*results.CreateProductResult, error)
	UpdateProduct(ctx context.Context, command *commands.UpdateProductCommand) (*results.UpdateProductResult, error)
	AddProductCategories(ctx context.Context, command *commands.AddProductCategoiesCommand) (*results.AddProductCategoriesResult, error)
	AddProductVariations(ctx context.Context, command *commands.AddProductVariationsCommand) (*results.AddProductVariationsResult, error)
}
