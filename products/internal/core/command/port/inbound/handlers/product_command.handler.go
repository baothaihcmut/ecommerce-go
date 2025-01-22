package handlers

import (
	"context"

	"github.com/baothaihcmut/Ecommerce-Go/products/internal/core/command/port/inbound/commands"
	"github.com/baothaihcmut/Ecommerce-Go/products/internal/core/command/port/inbound/results"
)

type ProductCommandHandler interface {
	CreateProduct(ctx context.Context, product *commands.CreateProductCommand) (*results.CreateProductResult, error)
}
