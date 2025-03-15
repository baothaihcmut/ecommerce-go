package valueobjects

import valueobjects "github.com/baothaihcmut/Ecommerce-go/products/internal/core/command/domain/common/value_objects"

type ProductImageId struct {
	ProductId ProductId
	Url       valueobjects.ImageLink
}

func NewProductImageId(productId ProductId, url valueobjects.ImageLink) ProductImageId {
	return ProductImageId{
		ProductId: productId,
		Url:       url,
	}
}
