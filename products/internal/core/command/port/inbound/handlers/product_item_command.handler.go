package handlers

import (
	"context"

	"github.com/baothaihcmut/Ecommerce-go/products/internal/core/command/port/inbound/commands"
	"github.com/baothaihcmut/Ecommerce-go/products/internal/core/command/port/inbound/results"
)

type ProductItemCommandHandler interface {
	CreateProductItem(ctx context.Context, command *commands.CreateProductItemCommand) (*results.CreateProductItemResult, error)
}
