package categories

import (
	valueobjects "github.com/baothaihcmut/Ecommerce-Go/products/internal/core/domain/aggregates/categories/value_objects"
)

type Category struct {
	Id               valueobjects.CategoryId
	Name             string
	ParentCategoryId []valueobjects.CategoryId
}

func NewCategory(id valueobjects.CategoryId, name string, parentCategoryId []valueobjects.CategoryId) *Category {
	return &Category{
		Id:               id,
		Name:             name,
		ParentCategoryId: parentCategoryId,
	}
}
