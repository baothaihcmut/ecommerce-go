package products

import (
	categoryValueobjects "github.com/baothaihcmut/Ecommerce-Go/products/internal/core/command/domain/aggregates/categories/value_objects"
	"github.com/baothaihcmut/Ecommerce-Go/products/internal/core/command/domain/aggregates/products/entities"
	valueobjects "github.com/baothaihcmut/Ecommerce-Go/products/internal/core/command/domain/aggregates/products/value_objects"
	"github.com/baothaihcmut/Ecommerce-Go/products/internal/core/command/domain/exceptions"
)

type Product struct {
	Id          valueobjects.ProductId
	Name        string
	Description string
	Unit        string
	ShopId      valueobjects.ShopId
	CategoryIds []categoryValueobjects.CategoryId
	Variations  []*entities.Variation
}

func checkVariationDuplicate(variations []string) error {
	variationSet := make(map[string]struct{})
	for _, variation := range variations {
		if _, exist := variationSet[variation]; !exist {
			variationSet[variation] = struct{}{}
		} else {
			return exceptions.ErrDuplicateVariation
		}
	}
	return nil
}
func checkCatgoryDuplicate(categoryIds []categoryValueobjects.CategoryId) error {
	categorySet := make(map[string]struct{})
	for _, categoryId := range categoryIds {
		if _, exist := categorySet[string(categoryId)]; !exist {
			categorySet[string(categoryId)] = struct{}{}
		} else {
			return exceptions.ErrCategoryExist
		}
	}
	return nil
}

func NewProduct(
	productId valueobjects.ProductId,
	name string,
	description string,
	unit string,
	shopId valueobjects.ShopId,
	categoryIds []categoryValueobjects.CategoryId,
	variations []string,
) (*Product, error) {
	err := checkCatgoryDuplicate(categoryIds)
	if err != nil {
		return nil, err
	}
	err = checkVariationDuplicate(variations)
	if err != nil {
		return nil, err
	}

	variationEntities := make([]*entities.Variation, len(variations))
	for idx, variation := range variations {
		variationEntities[idx] = entities.NewVariation(valueobjects.NewVariationId(productId, variation))
	}
	return &Product{
		Id:          productId,
		Name:        name,
		Description: description,
		Unit:        unit,
		ShopId:      shopId,
		CategoryIds: categoryIds,
		Variations:  variationEntities,
	}, nil
}
func (p *Product) AddVariation(variations []string) error {
	for _, variation := range variations {
		//check if variation exist in product
		for _, productVariation := range p.Variations {
			if variation == productVariation.Id.Name {
				return exceptions.ErrVariationExist
			}
		}
		p.Variations = append(p.Variations, entities.NewVariation(valueobjects.NewVariationId(p.Id, variation)))
	}
	return nil
}
func (p *Product) AddCategory(categoryId []categoryValueobjects.CategoryId) error {
	//check variation exist in product
	for _, val := range p.CategoryIds {
		for _, newCate := range categoryId {
			if val.IsEqual(newCate) {
				return exceptions.ErrCategoryExist
			}
		}
	}
	p.CategoryIds = append(p.CategoryIds, categoryId...)
	return nil
}

func (p *Product) CheckVariationBelongToProduct(variationId valueobjects.VariationId) bool {
	for _, productVariation := range p.Variations {
		if productVariation.Id.IsEqual(variationId) {
			return true
		}
	}
	return false
}
