package repositories

import (
	"context"

	"github.com/baothaihcmut/Ecommerce-go/libs/pkg/filter"
	"github.com/baothaihcmut/Ecommerce-go/libs/pkg/pagination"
	"github.com/baothaihcmut/Ecommerce-go/libs/pkg/sort"
	"github.com/baothaihcmut/Ecommerce-go/products/internal/core/query/projections/categories"
)

type CategoryQueryRepository interface {
	FindAllCategory(
		ctx context.Context,
		filters []filter.FilterParam,
		sorts []sort.SortParam,
		paginate pagination.PaginationParam,
	) (*pagination.PaginationResult[*categories.CategoryProjection], error)
	FindAllSubCategory(ctx context.Context, categoryId string) ([]*categories.CategoryProjection, error)
}
