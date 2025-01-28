package results

import productitems "github.com/baothaihcmut/Ecommerce-Go/products/internal/core/command/domain/aggregates/product_items"

type CreateProductItemResult struct {
	*productitems.ProductItem
}
