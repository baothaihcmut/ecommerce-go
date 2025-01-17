package handlers

import (
	"context"

	"github.com/baothaihcmut/Ecommerce-Go/products/internal/core/port/inbound/command/commands"
	"github.com/baothaihcmut/Ecommerce-Go/products/internal/core/port/inbound/command/results"
)

type CategoryCommandHandler interface {
	CreateCategory(context.Context, *commands.CreateCategoryCommand) (*results.CreateCategoryResult, error)
}
