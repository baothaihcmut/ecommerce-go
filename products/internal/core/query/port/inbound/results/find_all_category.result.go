package results

import (
	"github.com/baothaihcmut/Ecommerce-go/libs/pkg/pagination"
	"github.com/baothaihcmut/Ecommerce-go/products/internal/core/query/projections/categories"
)

type FindAllCategoryResult struct {
	*pagination.PaginationResult[*categories.CategoryProjection]
}
