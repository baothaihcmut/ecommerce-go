package results

import "github.com/baothaihcmut/Ecommerce-Go/products/internal/core/query/projections/categories"

type FindSubCategoriesResult struct {
	Categories []*categories.CategoryProjection
}
