package services

import (
	"context"

	"github.com/baothaihcmut/Ecommerce-Go/products/internal/core/query/port/inbound/handlers"
	"github.com/baothaihcmut/Ecommerce-Go/products/internal/core/query/port/inbound/queries"
	"github.com/baothaihcmut/Ecommerce-Go/products/internal/core/query/port/inbound/results"
	"github.com/baothaihcmut/Ecommerce-Go/products/internal/core/query/port/outbound/repositories"
)

type CategoryQueryService struct {
	repo repositories.CategoryQueryRepository
}

func NewCategoryQueryService(repo repositories.CategoryQueryRepository) handlers.CategoryQueryHandler {
	return &CategoryQueryService{
		repo: repo,
	}
}
func (c *CategoryQueryService) FindAllCategory(ctx context.Context, query *queries.FindAllCategoryQuery) (*results.FindAllCategoryResult, error) {
	res, err := c.repo.FindAllCategory(ctx, query.Filter, query.Sort, query.Pagination)
	if err != nil {
		return nil, err
	}
	return &results.FindAllCategoryResult{
		PaginationResult: res,
	}, nil
}
