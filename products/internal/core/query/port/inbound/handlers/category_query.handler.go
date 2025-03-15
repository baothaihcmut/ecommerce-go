package handlers

import (
	"context"

	"github.com/baothaihcmut/Ecommerce-go/products/internal/core/query/port/inbound/queries"
	"github.com/baothaihcmut/Ecommerce-go/products/internal/core/query/port/inbound/results"
)

type CategoryQueryHandler interface {
	FindAllCategory(context.Context, *queries.FindAllCategoryQuery) (*results.FindAllCategoryResult, error)
}
