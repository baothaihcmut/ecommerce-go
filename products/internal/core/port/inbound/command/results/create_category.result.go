package results

import valueobjects "github.com/baothaihcmut/Ecommerce-Go/products/internal/core/domain/aggregates/categories/value_objects"

type CreateCategoryResult struct {
	Id               valueobjects.CategoryId
	Name             string
	ParentCategoryId []valueobjects.CategoryId
}
