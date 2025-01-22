package results

import (
	"github.com/baothaihcmut/Ecommerce-Go/products/internal/core/command/domain/aggregates/categories"
)

type CreateCategoryResult struct {
	*categories.Category
}
