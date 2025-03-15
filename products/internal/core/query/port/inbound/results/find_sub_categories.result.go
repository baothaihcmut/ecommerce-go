package results

import "github.com/baothaihcmut/Ecommerce-go/products/internal/core/query/projections/categories"

type FindSubCategoriesResult struct {
	Categories []*categories.CategoryProjection
}
