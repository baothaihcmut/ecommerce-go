package results

import (
	"github.com/baothaihcmut/Ecommerce-go/products/internal/core/command/domain/aggregates/categories"
)

type CreateCategoryResult struct {
	*categories.Category
}
