package services

import (
	"context"

	"github.com/baothaihcmut/Ecommerce-Go/libs/pkg/utils"
	"github.com/baothaihcmut/Ecommerce-Go/products/internal/core/query/port/inbound/handlers"
	"github.com/baothaihcmut/Ecommerce-Go/products/internal/core/query/port/inbound/queries"
	"github.com/baothaihcmut/Ecommerce-Go/products/internal/core/query/port/inbound/results"
	"github.com/baothaihcmut/Ecommerce-Go/products/internal/core/query/port/outbound/repositories"
	"google.golang.org/protobuf/types/known/anypb"
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
	for idx, val := range query.Filter {
		anyVal := val.Value.(*anypb.Any)
		res, err := utils.ExtractGrpcAnyValue(anyVal)
		if err != nil {
			return nil, err
		}
		query.Filter[idx].Value = res
	}
	res, err := c.repo.FindAllCategory(ctx, query.Filter, query.Sort, query.Pagination)
	if err != nil {
		return nil, err
	}
	return &results.FindAllCategoryResult{
		PaginationResult: res,
	}, nil
}

func (c *CategoryQueryService) FindAllSubCategories(ctx context.Context, q *queries.FindSubCategoriesQuery) (*results.FindSubCategoriesResult, error) {
	res, err := c.repo.FindAllSubCategory(ctx, q.CategoryId)
	if err != nil {
		return nil, err
	}
	return &results.FindSubCategoriesResult{Categories: res}, nil
}
