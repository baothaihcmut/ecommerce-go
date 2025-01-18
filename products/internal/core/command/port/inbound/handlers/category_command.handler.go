package handlers

import (
	"context"

	"github.com/baothaihcmut/Ecommerce-Go/products/internal/core/command/port/inbound/commands"
	"github.com/baothaihcmut/Ecommerce-Go/products/internal/core/command/port/inbound/results"
)

type CategoryCommandHandler interface {
	CreateCategory(context.Context, *commands.CreateCategoryCommand) (*results.CreateCategoryResult, error)
}
