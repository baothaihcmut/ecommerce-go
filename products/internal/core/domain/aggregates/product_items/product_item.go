package productitems

import (
	"errors"

	"github.com/baothaihcmut/Ecommerce-Go/products/internal/core/domain/aggregates/product_items/entities"
	valueobjects "github.com/baothaihcmut/Ecommerce-Go/products/internal/core/domain/aggregates/product_items/value_objects"
	productValueobjects "github.com/baothaihcmut/Ecommerce-Go/products/internal/core/domain/aggregates/products/value_objects"
)

var (
	ErrDuplicateVariation     = errors.New("duplicate variation in one product")
	ErrMismatchVariationValue = errors.New("you cannot add this varition value to this item")
)

type ProductItem struct {
	Id              valueobjects.ProductItemId
	Sku             valueobjects.SKU
	Price           valueobjects.ProductPrice
	Quantity        valueobjects.ProductQuantity
	ProductId       productValueobjects.ProductId
	VariationValues []*entities.VariationValue
}

func NewProductItem(
	id valueobjects.ProductItemId,
	sku valueobjects.SKU,
	price valueobjects.ProductPrice,
	quantity valueobjects.ProductQuantity,
	productId productValueobjects.ProductId,
	variationValues []*entities.VariationValue,
) (*ProductItem, error) {
	//check if product item have two same variationvalues of 1 variation
	mapVariation := make(map[string]bool, len(variationValues))
	for _, variation := range variationValues {
		//check match product id
		if !variation.VariationId.ProductId.IsEqual(productId) {
			return nil, ErrMismatchVariationValue
		}
		key := string(variation.VariationId.ProductId) + string(variation.VariationId.Name)
		if _, exist := mapVariation[key]; exist {
			return nil, ErrDuplicateVariation
		} else {
			mapVariation[key] = true
		}
	}

	return &ProductItem{
		Id:              id,
		Sku:             sku,
		Price:           price,
		Quantity:        quantity,
		ProductId:       productId,
		VariationValues: variationValues,
	}, nil
}
func (p *ProductItem) AddVariationValues(variationValues []*entities.VariationValue) error {
	for _, variationValue := range p.VariationValues {
		for _, newVariation := range variationValues {
			//check match product id
			if !variationValue.VariationId.ProductId.IsEqual(p.ProductId) {
				return ErrMismatchVariationValue
			}
			if variationValue.VariationId.IsEqual(newVariation.VariationId) {
				return ErrDuplicateVariation
			}
		}
	}
	p.VariationValues = append(p.VariationValues, variationValues...)
	return nil
}

func (p *ProductItem) IncreProductItemQuantity(quantity valueobjects.ProductQuantity) {
	p.Quantity = p.Quantity.IncreQuantity(quantity)
}

func (p *ProductItem) DecreProductItemQuantity(quantity valueobjects.ProductQuantity) error {
	newQuantity, err := p.Quantity.DecreQuantity(quantity)
	if err != nil {
		return err
	}
	p.Quantity = newQuantity
	return nil
}
