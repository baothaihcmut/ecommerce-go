package products

import (
	"errors"

	categoryValueobjects "github.com/baothaihcmut/Ecommerce-Go/products/internal/core/domain/aggregates/categories/value_objects"
	"github.com/baothaihcmut/Ecommerce-Go/products/internal/core/domain/aggregates/products/entities"
	valueobjects "github.com/baothaihcmut/Ecommerce-Go/products/internal/core/domain/aggregates/products/value_objects"
	shopValueobjects "github.com/baothaihcmut/Ecommerce-Go/products/internal/core/domain/aggregates/shops/value_objects"
)

var (
	ErrVariationExist              = errors.New("variation exist in product")
	ErrCategoryExist               = errors.New("category exist in product")
	ErrVariationNotBelongToProduct = errors.New("variation not belong to product")
)

type Product struct {
	Id          valueobjects.ProductId
	Name        string
	Description string
	Unit        string
	ShopId      shopValueobjects.ShopId
	CategoryIds []categoryValueobjects.CategoryId
	Variations  []*entities.Variation
}

func NewProduct(
	id valueobjects.ProductId,
	name string,
	description string,
	unit string,
	shopId shopValueobjects.ShopId,
	categoryIds []categoryValueobjects.CategoryId,
	variations []*entities.Variation,
) *Product {

	return &Product{
		Id:          id,
		Name:        name,
		Description: description,
		Unit:        unit,
		ShopId:      shopId,
		CategoryIds: categoryIds,
		Variations:  variations,
	}
}
func (p *Product) AddVariation(variations []*entities.Variation) error {
	for _, variation := range variations {
		//check if variation belong to product
		if !variation.Id.ProductId.IsEqual(p.Id) {
			return ErrVariationNotBelongToProduct
		}
		//check if variation exist in product
		for _, productVariation := range p.Variations {
			if !productVariation.Id.IsEqual(variation.Id) {
				return ErrVariationExist
			}
		}
		p.Variations = append(p.Variations, variation)
	}
	return nil
}
func (p *Product) AddCategory(categoryId []categoryValueobjects.CategoryId) error {
	//check variation exist in product
	for _, val := range p.CategoryIds {
		for _, newCate := range categoryId {
			if val.IsEqual(newCate) {
				return ErrCategoryExist
			}
		}
	}
	p.CategoryIds = append(p.CategoryIds, categoryId...)
	return nil
}
