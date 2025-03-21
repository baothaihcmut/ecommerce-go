package handlers

import (
	"context"

	"github.com/baothaihcmut/Ecommerce-go/products/internal/core/port/inbound/commands"
	"github.com/baothaihcmut/Ecommerce-go/products/internal/core/port/inbound/results"
)

type ProductHandler interface {
	CreateProduct(context.Context, *commands.CreateProductCommand) (*results.CreateProductResult, error)
}
