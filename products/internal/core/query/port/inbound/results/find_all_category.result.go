package results

import (
	"github.com/baothaihcmut/Ecommerce-Go/libs/pkg/pagination"
	"github.com/baothaihcmut/Ecommerce-Go/products/internal/core/query/projections/categories"
)

type FindAllCategoryResult struct {
	*pagination.PaginationResult[*categories.CategoryProjection]
}
